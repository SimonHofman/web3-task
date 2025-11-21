// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {MyNftAuction} from "./MyNftAuction.sol";

contract MyNftAuctionStore {
    address[] public auctions;

    mapping(uint256 tokenId => MyNftAuction) public auctionMap;

    event AuctionCreated(address indexed auctionAddress, uint256 tokenId);

    function createAuction(
        uint256 duration,
        uint256 startPrice,
        address nftContractAddress,
        uint256 tokenId
    ) external returns (address){
        MyNftAuction auction = new MyNftAuction();
        auction.initialize();
        auction.createAuction(duration, startPrice, nftContractAddress, tokenId);
        auctions.push(address(auction));
        auctionMap[tokenId] = auction;

        emit AuctionCreated(address(auction), tokenId);
        return address(auction);
    }

    function getAuctions() external view returns (address[] memory) {
        return auctions;
    }

    function getAuction(uint256 tokenId) external view returns (address) {
        require(tokenId < auctions.length, "tokenId out of bounds");
        return auctions[tokenId];
    }
}
