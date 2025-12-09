//SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract NFT is ERC721,ERC721URIStorage,Ownable {
    uint256 private _tokenIdCounter;

    constructor() ERC721("AuctionNFT","ANFT") Ownable(msg.sender) {}

    function mintNFT(address to ,string memory uri) public onlyOwner returns(uint256){
        uint256 tokenId = _tokenIdCounter++;
        _safeMint(to,tokenId);
        _setTokenURI(tokenId,uri);
        return tokenId;
    }
    function burn(uint256 tokenId) public virtual onlyOwner {
        _burn(tokenId);
    }
    function supportsInterface(bytes4 interfaceId) public view override(ERC721,ERC721URIStorage) returns(bool){
        return super.supportsInterface(interfaceId);
    }
    function tokenURI(uint256 tokenId) 
        public 
        view 
        override(ERC721, ERC721URIStorage) 
        returns (string memory) 
    {
        return super.tokenURI(tokenId);
    }
}