// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
import "./MyNftAuction.sol";

contract MyNftAuctionV2 is MyNftAuction {
    function testHello() public pure returns (string memory) {
        return "hello world!";
    }
}
