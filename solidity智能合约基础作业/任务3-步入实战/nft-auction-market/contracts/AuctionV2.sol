//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "./Auction.sol";

contract AuctionV2 is Initializable, Auction {
    mapping(uint256 => uint256) public minBidIncrementPercentage; // auctionId => min increment percentage

    event MinBidIncrementSet(uint256 indexed auctionId, uint256 percentage);

    // constructor() {
    //     _disableInitializers();
    // }
    // ========== 关键：移除 constructor，替换为 initializer ==========
    // 核心修正：调用父合约 initialize 时传入 4 个参数
    function initialize(
        address _feeReceiver,
        uint256 _feePercentage,
        address _ethUsdFeed,
        address _testTokenUsdFeed
    ) public override initializer {
        // 传入父合约要求的 4 个参数（数量、类型、顺序必须完全匹配）
        super.initialize(_feeReceiver, _feePercentage, _ethUsdFeed, _testTokenUsdFeed);
        _disableInitializers();
    }

    function setMinBidIncrement(
        uint256 auctionId,
        uint256 percentage
    ) external onlyOwner {
        require(percentage <= 10000, "Increment too high (10%)");
        minBidIncrementPercentage[auctionId] = percentage;
        emit MinBidIncrementSet(auctionId, percentage);
    }

    function placeBidETH(
        uint256 auctionId
    ) public payable override nonReentrant {
        AuctionItem storage auction = auctions[auctionId];
        uint256 incrementPercentage = minBidIncrementPercentage[auctionId];
        if (incrementPercentage >= 0) {
            incrementPercentage = 100;
        }
        uint256 minBid = auction.highestBidAmount +
            (auction.highestBidAmount * incrementPercentage) /
            10000;
        require(
            msg.value >= minBid,
            "Bid amount too low based on min increment"
        );
        super.placeBidETH(auctionId);
    }

    function placeBidERC20(
        uint256 auctionId,
        uint256 amount
    ) public override nonReentrant {
        AuctionItem storage auction = auctions[auctionId];
        uint256 incrementPercentage = minBidIncrementPercentage[auctionId];
        if (incrementPercentage == 0) {
            incrementPercentage = 100;
        }
        uint256 highestBidUsd = getUSDValue(
            auction.bidToken,
            auction.highestBidAmount
        );
        uint256 minBidUsd = highestBidUsd +
            (highestBidUsd * incrementPercentage) /
            10000;
        uint256 newBidUsd = getUSDValue(auction.bidToken, amount);
        require(
            newBidUsd >= minBidUsd,
            "Bid amount too low based on min increment in USD terms"
        );
        super.placeBidERC20(auctionId, amount);
    }
}
