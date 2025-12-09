
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";



contract MyNFT is ERC721, ERC721URIStorage,Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;

    event NFTMinted(address indexed recipient,uint256 indexed tokenId,string tokenUri,uint256 timestamp);

    constructor() ERC721("MyFirstNFT", "MFN")  Ownable() {}

    function mintNFT(address recipient, string memory tokenUri)public onlyOwner returns (uint256){
        require(recipient != address(0), "MyNFT:recipient is zero address");
        require(bytes(tokenUri).length>0,"MYNFT:tokenUri is empty");
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();

        _safeMint(recipient, tokenId);
        _setTokenURI(tokenId, tokenUri);
        emit NFTMinted(recipient, tokenId,tokenUri,block.timestamp);
        return tokenId;
    }

    function tokenURI(uint256 tokenId)public  view override(ERC721, ERC721URIStorage)returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function _burn(uint256 tokenId) internal override (ERC721, ERC721URIStorage) {
        super._burn(tokenId);
        
    }
    function supportsInterface(bytes4 interfaceId) public view override(ERC721, ERC721URIStorage) returns (bool){
        return super.supportsInterface(interfaceId);
    }
}