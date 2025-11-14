// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BeggingContract {
    mapping(address => uint256) public donates;
    uint256 private totalDonated;
    address private owner;

    constructor() {
        owner = msg.sender;
    }

    function donate() public payable {
        donates[msg.sender] += msg.value;
        totalDonated += msg.value;
    }

    function getDonation(address from) public view returns(uint256) {
        return donates[from];
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "must be owner");
        _;
    }

    function withdraw() public onlyOwner {
        payable(msg.sender).transfer(address(this).balance);
    }
}
