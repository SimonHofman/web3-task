// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


contract MyNFT1 is ERC721URIStorage, Ownable{
    uint256 private _tokenIds;  // 直接使用 uint256 替代 Counters.Counter

    constructor(address initialOwer)
        ERC721("MYNFT1", "MYNFT1")
        Ownable(initialOwer)
    {}

    function mintNFT(address to, string memory tokenURI) public onlyOwner {
        _tokenIds++;  // 直接递增
        uint256 tokenId = _tokenIds;
        // 铸造
        _safeMint(to, tokenId);
        // 设置NFT的元数据
        _setTokenURI(tokenId, tokenURI);
    }
  
}
