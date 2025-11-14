// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BeggingContract2 {
    mapping(address => uint256) private _donates;
    address private _owner;

    event Donation(address indexed donate, uint256 amount);
    event Withdraw(address indexed owner, uint256 amcount);
    event FallbackCalled(address sender, uint256 amount, bytes data);
    event OwnershipTransferred(
        address indexed previousOwner,
        address indexed newOwner
    );

    receive() external payable {
        donate();
    }

    fallback() external payable {
        require(
            msg.value > 0,
            "BeggingContract2: The donation amount must be greater than 0"
        );
        _donates[msg.sender] += msg.value;
        emit FallbackCalled(msg.sender, msg.value, msg.data);
    }

    constructor() {
        _owner = msg.sender;
        emit OwnershipTransferred(address(0), _owner);
    }

    function donate() public payable {
        require(
            msg.value > 0,
            "BeggingContract2: The donation amount must be greater than 0"
        );
        _donates[msg.sender] += msg.value;
        emit Donation(msg.sender, msg.value);
    }

    function getDonation(address from) public view returns(uint256) {
        return _donates[from];
    }

    modifier onlyOwner() {
        require(msg.sender == _owner, "must be owner");
        _;
    }

    function withdraw() public onlyOwner {
        // 获取当前合约的总金额
        uint256 amount = address(this).balance;
        require(amount > 0, "BeggingContract2: No account can be withdraw");
        payable(msg.sender).transfer(address(this).balance);
        emit Withdraw(msg.sender, amount);
    }

    // 转移所有者权限
    function transferOwnership(address newOwner) public onlyOwner {
        require(
            newOwner != address(0),
            "BeggingContract2: new owner is the zero address"
        );
        _owner = newOwner;
        emit OwnershipTransferred(_owner, newOwner);
    }

    function getOwner() public view returns(address) {
        return _owner;
    }
}
