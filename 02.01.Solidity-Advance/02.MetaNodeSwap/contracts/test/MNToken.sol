// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MNToken is ERC20 {
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {}

    function mint(address recipient, uint256 quantity) public payable {
        require(quantity > 0, "quantity must be greater than 0");
        require(quantity < 10_000_000 * 10**18, "quantity must be less then 10000000 * 10 ** 18");
        _mint(recipient, quantity);
    }
}
