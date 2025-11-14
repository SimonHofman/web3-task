// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";


contract MyNFT1 is ERC721URIStorage, Ownable{
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor(address initialOwer)
        ERC721("MYNFT1", "MYNFT1")
        Ownable(initialOwer)
    {}

    function mintNFT(address to, string memory tokenURI) public onlyOwner {
        _tokenIds.increment();
        uint256 tokenId = _tokenIds.current();
        // 铸造
        _safeMint(to, tokenId);
        // 设置NFT的元数据
        _setTokenURI(tokenId, tokenURI);
    }
  
}
