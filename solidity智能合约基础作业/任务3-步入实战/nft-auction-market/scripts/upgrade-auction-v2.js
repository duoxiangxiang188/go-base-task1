const { ethers, upgrades } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();

    console.log(`Upgrading Auction Contract with account: ${deployer.address}`);

    const PROXY_ADDRESS = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"; // Replace with the deployed proxy address of the Auction contract

    const AuctionV2 = await ethers.getContractFactory("AuctionV2");
    const upgraded = await upgrades.upgradeProxy(PROXY_ADDRESS, AuctionV2);
    console.log("Auction Contract upgraded to V2 at address:", await upgraded.getAddress());
    console.log(`Aution V2 implementation address: ${await upgrades.erc1967.getImplementationAddress(await upgraded.getAddress())}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exitCode = 1;
    });