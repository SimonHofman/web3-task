// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract count {
    uint256 private num = 0;

    event Counted(uint256 num);

    function get() public view returns (uint256) {
        return num;
    }

    function inc() public {
        num++;
        emit Counted(num);
    }
}
