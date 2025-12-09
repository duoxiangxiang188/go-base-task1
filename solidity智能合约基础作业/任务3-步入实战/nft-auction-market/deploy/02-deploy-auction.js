const { ethers, upgrades } = require("hardhat");

module.exports = async ({ getNamedAccounts, deployments }) => {
    const { deploy, log } = deployments;
    const { deployer } = await getNamedAccounts();

    const ETH_USD_FEED = "0x0000000000000000000000000000000000000000"; // Replace with actual ETH/USD price feed address
    const TEST_TOKEN_USD_FEED = "0x0000000000000000000000000000000000000000"; // Replace with actual Test Token/USD price feed address

    log("----------------------------------------------------");
    log("Deploying Auction Contract with UUPS Proxy...");

    const Auction = await ethers.getContractFactory("Auction");
    const auctionProxy = await upgrades.deployProxy(
        Auction,
        [deployer, 100, ETH_USD_FEED, TEST_TOKEN_USD_FEED],
        {
            initializer: 'initialize',
            kind: 'uups',
           
        }
    );

    await auctionProxy.waitForDeployment();
    const auctionProxyAddress = await auctionProxy.getAddress();
    const auctionImplAddress = await upgrades.erc1967.getImplementationAddress(auctionProxyAddress);
    log("Auction Contract deployed to:", auctionProxyAddress);
    log(`Aution implementation address: ${auctionImplAddress}`);
};
module.exports.tags = ["Auction", "UUPS", "all"];