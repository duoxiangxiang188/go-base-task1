const {expect} = require("chai");
const {ethers} = require("hardhat");

describe("NFT Contract", function () {
    let NFT;
    let nft;
    let owner;
    let addr1;
    this.beforeEach(async function () {
        NFT = await ethers.getContractFactory("NFT");
        [owner, addr1] = await ethers.getSigners();
        nft = await NFT.deploy();
        await nft.waitForDeployment();
    });

    it("Should mint a new NFT and assign it to the caller", async function () {
        const uri ="ipfs://QmXYZ...";
        await nft.mintNFT(addr1.address, uri);
        expect(await nft.balanceOf(addr1.address)).to.equal(1);
        expect(await nft.tokenURI(0)).to.equal(uri);
        expect(await nft.ownerOf(0)).to.equal(addr1.address);
    });

    it("Should not allow non-owner to mint", async function () {
        const uri ="ipfs://QmXYZ...";
        await expect(nft.connect(addr1).mintNFT(addr1.address, uri)).to.be.revertedWithCustomError(nft, "OwnableUnauthorizedAccount");
    });
});