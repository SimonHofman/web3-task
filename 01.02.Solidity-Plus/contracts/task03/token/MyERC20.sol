// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

contract MyERC20  is ERC20, ERC20Permit{
    constructor() ERC20("MyToken", "MyERC20") ERC20Permit("MyToken"){
        _mint(msg.sender, 100000 * 10 ** decimals());
    }
}
