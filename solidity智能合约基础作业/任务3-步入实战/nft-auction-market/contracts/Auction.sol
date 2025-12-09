// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
// 导入OpenZeppelin升级合约库
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

// 导入ERC20相关接口与工具
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract Auction is Initializable ,OwnableUpgradeable,UUPSUpgradeable,ReentrancyGuardUpgradeable,IERC721Receiver {

    using SafeERC20 for IERC20;
    enum AuctionStatus {Created,Active,Ended }

    struct AuctionItem {
        uint256 auctionId;
        address nftContract;
        uint256 tokenId;
        address seller;
        uint256 startTime;
        uint256 endTime;
        address highestBidder;
        uint256 highestBidAmount;
        address bidToken;
        AuctionStatus status;
        uint256 feePercentage;
        
    }

    mapping (address => AggregatorV3Interface) public priceFeeds;
    mapping(address => uint256) public testUsdRate;
    address public _testToken;

    uint256 private _auctionIdCounter;
    mapping (uint256 => AuctionItem) public auctions;
    mapping (uint256 => mapping(address => uint256)) public bids; // auctionId => bidder => amount

    event AuctionCreated(uint256 indexed auctionId, address seller,address indexed nftContract, uint256 indexed tokenId);
    event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount,address bidToken);
    event AuctionEnded(uint256 indexed auctionId, address indexed winner, uint256 amount);
    event FeeUpdated(uint256 newFeePercentage);
    event PriceFeedUpdated(address indexed token, address indexed feed);



    uint256 public FeePencentage; // e.g., 250 for 2.5%
    address public feeReceiver;

    // constructor(address testTokenAddr) {
    //     _testToken = testTokenAddr;
    //     _disableInitializers();
    // }   
    function setPriceFeed(address token,address feed) external onlyOwner {
        priceFeeds[token] = AggregatorV3Interface(feed);
    }

    function setTestUsdRate(address token,uint256 rate) external onlyOwner {
        testUsdRate[token] = rate;
    }
    function initialize(address _feeReceiver, uint256 _feePercentage, address _ethUsdFeed,address _testTokenUsdFeed) public virtual  initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
        feeReceiver = _feeReceiver;
        FeePencentage = _feePercentage;
        priceFeeds[address(0)] = AggregatorV3Interface(_ethUsdFeed); // ETH
        priceFeeds[_testToken] = AggregatorV3Interface(_testTokenUsdFeed); // TestToken
    }
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
    function onERC721Received(
        address, // operator（无需使用）
        address, // from（无需使用）
        uint256, // tokenId（无需使用）
        bytes calldata // data（无需使用）
    ) external pure override returns (bytes4) {
        // 返回接口的选择器，表示接收成功
        return IERC721Receiver.onERC721Received.selector;
    }

    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startTime,
        uint256 endTime,
        address bidToken,
        uint256 initialBid
    ) external {
        require(IERC721(nftContract).ownerOf(tokenId)== msg.sender,"You are not the owner of this NFT");
        require(startTime > block.timestamp,"Start time should be in the future");
        require(endTime >startTime,"End time should be after start time");
        IERC721(nftContract).safeTransferFrom(msg.sender,address(this),tokenId);

        _auctionIdCounter++;
        uint256 auctionId = _auctionIdCounter;
        uint256 dynamicFee = initialBid >1 ether ? 50 : 100; // 0.5% for bids > 1 ETH, else 1%
        auctions[auctionId] = AuctionItem({
            auctionId: auctionId,
            nftContract: nftContract,
            tokenId: tokenId,
            seller: msg.sender,
            startTime: startTime,
            endTime: endTime,
            highestBidder: address(0),
            highestBidAmount: initialBid,
            bidToken: bidToken,
            status: AuctionStatus.Created,
            feePercentage:  dynamicFee
        });
        emit AuctionCreated(auctionId,msg.sender,nftContract,tokenId);

    }

    function startAuction(uint256 auctionId) external onlyOwner {
        AuctionItem storage auction = auctions[auctionId];
        require(auction.status == AuctionStatus.Created,"Auction is not in Created status");
        require(block.timestamp >= auction.startTime,"Auction start time not reached");
        auction.status = AuctionStatus.Active;
    }
        // 出价（ETH）
    function placeBidETH(uint256 auctionId)public payable nonReentrant virtual{
        AuctionItem storage auction = auctions[auctionId];
        require(auction.status == AuctionStatus.Active,"Auction is not active");
        require(block.timestamp >= auction.startTime && block.timestamp <= auction.endTime,"Auction is not active");
        require(auction.bidToken == address(0),"This auction does not accept ETH bids");
        uint256 minBidAmount = auction.highestBidAmount + (auction.highestBidAmount * auction.feePercentage) / 10000;
        require(msg.value >= minBidAmount,"Bid amount too low");

        // Refund previous highest bidder
        if (auction.highestBidder != address(0)) {
            payable(auction.highestBidder).transfer(auction.highestBidAmount);
        }

        auction.highestBidder = msg.sender;
        auction.highestBidAmount = msg.value;
        bids[auctionId][msg.sender] += msg.value;

        emit BidPlaced(auctionId,msg.sender,msg.value,address(0));
    }
    function placeBidERC20(uint256 auctionId,uint256 amount) public nonReentrant virtual{
        AuctionItem storage auction = auctions[auctionId];
        require(auction.status == AuctionStatus.Active,"Auction is not active");
        require(block.timestamp >= auction.startTime && block.timestamp <= auction.endTime,"Auction is not active");
        require(auction.bidToken != address(0),"This auction only accepts ETH bids");
        
        uint256 newBidUsd = getUSDValue(auction.bidToken,amount);
        uint256 highestBidUsd = getUSDValue(auction.bidToken,auction.highestBidAmount);
        uint256 minBidUsd = highestBidUsd + (highestBidUsd * auction.feePercentage) / 10000;
        require(newBidUsd >= minBidUsd,"Bid amount too low in USD terms");


        IERC20(auction.bidToken).safeTransferFrom(msg.sender,address(this),amount);

        // Refund previous highest bidder
        if (auction.highestBidder != address(0)) {
            IERC20(auction.bidToken).safeTransfer(auction.highestBidder,auction.highestBidAmount);
        }

        auction.highestBidder = msg.sender;
        auction.highestBidAmount = amount;
        bids[auctionId][msg.sender] += amount;

        emit BidPlaced(auctionId,msg.sender,amount,auction.bidToken);
    }
    function endAuction(uint256 auctionId) external nonReentrant {
        AuctionItem storage auction = auctions[auctionId];
        require(auction.status == AuctionStatus.Active,"Auction is not active");
        require(block.timestamp > auction.endTime,"Auction has not ended yet");

        auction.status = AuctionStatus.Ended;

        if (auction.highestBidder != address(0)) {
            // Transfer NFT to highest bidder
            IERC721(auction.nftContract).safeTransferFrom(address(this),auction.highestBidder,auction.tokenId);

            // Calculate fee and transfer funds
            uint256 feeAmount = (auction.highestBidAmount * FeePencentage) / 10000;
            uint256 sellerAmount = auction.highestBidAmount - feeAmount;

            if (auction.bidToken == address(0)) {
                // ETH
                payable(feeReceiver).transfer(feeAmount);
                payable(auction.seller).transfer(sellerAmount);
            } else {
                // ERC20
                IERC20(auction.bidToken).safeTransfer(feeReceiver,feeAmount);
                IERC20(auction.bidToken).safeTransfer(auction.seller,sellerAmount);
            }
        } else {
            // No bids received, return NFT to seller
            IERC721(auction.nftContract).safeTransferFrom(address(this),auction.seller,auction.tokenId);
        }

        emit AuctionEnded(auctionId,auction.highestBidder,auction.highestBidAmount);
    }
    function getUSDValue(address token,uint256 amount) public view returns(uint256){
        if (block.chainid == 31337) {
            uint256 rate = testUsdRate[token];
            require(rate > 0, "Test USD rate not set for this token");
            return (rate * amount) / 1e18;
            
        }


        AggregatorV3Interface priceFeed = priceFeeds[token];
        require(address(priceFeed) != address(0),"Price feed not available for this token");
        (,int256 price,,uint256 updatedAt,) = priceFeed.latestRoundData();
        require(price > 0, "Invalid price"); // 确保价格为正
        require(block.timestamp - updatedAt < 3600, "Price feed outdated"); // 价格更新时间不超过1小时
        uint8 decimals = priceFeed.decimals();
        return (uint256(price) * amount) / (10 ** decimals);
    }
    function updatePriceFeed(address token,address feed) external onlyOwner {
        priceFeeds[token] = AggregatorV3Interface(feed);
        emit PriceFeedUpdated(token,feed);
    }
    function updateFeeReceiver(address _feeReceiver) external onlyOwner {
        feeReceiver = _feeReceiver;
    }

    function updateFeePercentage(uint256 _feePercentage) external onlyOwner {
        require(_feePercentage <= 1000,"Fee percentage too high"); // Max 10%
        FeePencentage = _feePercentage;
        emit FeeUpdated(_feePercentage);
    }
}