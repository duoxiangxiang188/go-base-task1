const { expect } = require('chai');
const { ethers } = require('hardhat');
const { time } = require('@nomicfoundation/hardhat-network-helpers');


describe('Auction Contract', function () {
    let NFT;
    let nft;
    let TestERC20;
    let testERC20;
    let Auction;
    let auction;
    let owner;
    let seller;
    let bidder1;
    let bidder2;
    let feeReceiver;

    const ETH_USD_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306"; // Sepolia ETH/USD
    const TEST_TOKEN_USD_FEED = "0x14866185B1962B63C3Ea9E03Bc1da838bab34C19"; // Sepolia LINK/USD

    this.beforeEach(async function () {
        [owner, seller, bidder1, bidder2, feeReceiver] = await ethers.getSigners();
        NFT = await ethers.getContractFactory('NFT');
        nft = await NFT.deploy();
        await nft.waitForDeployment();

        TestERC20 = await ethers.getContractFactory('TestERC20');
        testERC20 = await TestERC20.deploy();
        await testERC20.waitForDeployment();
        Auction = await ethers.getContractFactory('Auction');
        auction = await upgrades.deployProxy(
            Auction,
            [feeReceiver.address, 100, ETH_USD_FEED, TEST_TOKEN_USD_FEED],
            { initializer: 'initialize', kind: 'uups' }
        );
        await auction.waitForDeployment();

        await nft.mintNFT(seller.address, "ipfs://QmSellerNFT");
        await nft.connect(seller).approve(await auction.getAddress(), 0);
        await testERC20.mint(bidder1.address, ethers.parseUnits("1000", 18));
        await testERC20.mint(bidder2.address, ethers.parseUnits("1000", 18));
    });
    it('Should create an auction', async function () {
        const startTime = (await time.latest()) + 60; // 1 minute from now
        const endTime = startTime + 3600; // 1 hour duration

        await expect(auction.connect(seller).createAuction(
            await nft.getAddress(),
            0,
            startTime,
            endTime,
            ethers.ZeroAddress,
            ethers.parseUnits("0.1", 18)
        )).to.emit(auction, 'AuctionCreated').withArgs(1, seller.address, await nft.getAddress(), 0);
        const auctionItem = await auction.auctions(1);
        expect(auctionItem.seller).to.equal(seller.address);
        expect(auctionItem.status).to.equal(0); // Created
    });

    it('Should place bids and end auction', async function () {
        const startTime = (await time.latest()) + 60; // 1 minute from now
        const endTime = startTime + 3600; // 1 hour duration
        // 步骤1：竞拍前记录 feeReceiver 的初始余额
        const initialFeeReceiverBalance = await ethers.provider.getBalance(feeReceiver.address);

        await auction.connect(seller).createAuction(
            await nft.getAddress(),
            0,
            startTime,
            endTime,
            ethers.ZeroAddress,
            ethers.parseUnits("0.1", 18)
        );

        await time.increaseTo(startTime + 10); // Move time to after auction start
        await auction.startAuction(1);

        await expect(auction.connect(bidder1).placeBidETH(1, { value: ethers.parseUnits("2", 18) }))
            .to.emit(auction, 'BidPlaced').withArgs(1, bidder1.address, ethers.parseUnits("2", 18), ethers.ZeroAddress);
        await expect(auction.connect(bidder2).placeBidETH(1, { value: ethers.parseUnits("3", 18) }))
            .to.emit(auction, 'BidPlaced').withArgs(1, bidder2.address, ethers.parseUnits("3", 18), ethers.ZeroAddress);


        await time.increaseTo(endTime + 10); // Move time to after auction end
        await expect(auction.endAuction(1))
            .to.emit(auction, 'AuctionEnded').withArgs(1, bidder2.address, ethers.parseUnits("3", 18));


        expect(await nft.ownerOf(0)).to.equal(bidder2.address);
        const fee = ethers.parseUnits("3", 18) * 100n / 10000n; // 1% fee
        expect(await ethers.provider.getBalance(feeReceiver.address) - initialFeeReceiverBalance).to.equal(fee);
        expect(await ethers.provider.getBalance(seller.address)).to.be.gt(ethers.parseUnits("97", 18));
    });

    it('Should place ERC20  bids with USD comparison', async function () {
        // 2. 部署模拟PriceFeed（1 TUSD = 1 USD，价格=1e8，8位小数）
        const MockPriceFeed = await ethers.getContractFactory("MockPriceFeed");
        const mockPriceFeed = await MockPriceFeed.deploy(
            ethers.parseUnits("1", 8), // price=1e8（1 USD/TUSD）
            8 // decimals=8（对齐Chainlink）
        );
        await mockPriceFeed.waitForDeployment();
        const mockPriceFeedAddr = await mockPriceFeed.getAddress();
        console.log("MockPriceFeed deployed to:", mockPriceFeedAddr);

        await auction.setPriceFeed(await testERC20.getAddress(), mockPriceFeedAddr);
        console.log("PriceFeed mapped: TestERC20 → MockPriceFeed");


        const startTime = (await time.latest()) + 60; // 1 minute from now
        const endTime = startTime + 3600; // 1 hour duration

        await auction.connect(seller).createAuction(
            await nft.getAddress(),
            0,
            startTime,
            endTime,
            await testERC20.getAddress(), // Test ERC20
            ethers.parseUnits("100", 18)
        );
        await time.increaseTo(startTime + 10); // Move time to after auction start
        await auction.startAuction(1);
        const testUsdRate = ethers.parseUnits("1", 8); // 假设1 TestToken = 1 USD
        await auction.connect(owner).setTestUsdRate(await testERC20.getAddress(), testUsdRate);

        // ========== 核心补充：竞价前授权 ==========
        const bidAmount1 = ethers.parseUnits("200", 18);
        // bidder1 授权拍卖合约花费其 TestERC20
        await testERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount1);
        // bidder2 同理，先授权再出价
        const bidAmount2 = ethers.parseUnits("300", 18);
        await testERC20.connect(bidder2).approve(await auction.getAddress(), bidAmount2);

        
        await auction.connect(bidder1).placeBidERC20(1, ethers.parseUnits("200", 18));
        await expect(auction.connect(bidder2).placeBidERC20(1, ethers.parseUnits("300", 18)))
            .to.emit(auction, 'BidPlaced').withArgs(1, bidder2.address, ethers.parseUnits("300", 18), await testERC20.getAddress());

        await time.increaseTo(endTime + 10); // Move time to after auction end
        await auction.endAuction(1);

        expect(await nft.ownerOf(0)).to.equal(bidder2.address);

    });
});
