// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BeggingContractExtra {
    // 时间限制（开始时间和结束时间）
    uint256 public _startTime;
    uint256 public _endTime;
    mapping(address => uint256) private _donors;
    address private _owner;
    // 捐赠金额前三
    address[3] private _top3Donors;

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
        _donors[msg.sender] += msg.value;
        emit FallbackCalled(msg.sender, msg.value, msg.data);
    }

    constructor(uint256 startTime, uint256 endTime) {
        _startTime = startTime;
        _endTime = endTime;
        _owner = msg.sender;
        _top3Donors = [address(0), address(0), address(0)];
        emit OwnershipTransferred(address(0), _owner);
    }

    modifier donationPeriod() {
        require(
            block.timestamp >= _startTime && block.timestamp <= _endTime,
            "BeggingContractExtra: Donations are only allowed during the specified period"
        );
        _;
    }

    function donate() public payable donationPeriod {
        require(
            msg.value > 0,
            "BeggingContract2: The donation amount must be greater than 0"
        );
        _donors[msg.sender] += msg.value;
        _updateTop3Donors();
        emit Donation(msg.sender, msg.value);
    }

    function _updateTop3Donors() private {
        // 检查是否已在排行榜中
        for (uint i = 0; i < 3; i++) {
            if (_top3Donors[i] == msg.sender) {
                // 触发内部排序
                for(uint j = i; j > 0; j--) {
                    address tAddress;
                    if (_donors[_top3Donors[j]] > _donors[_top3Donors[j-1]]) {
                        tAddress =  _top3Donors[j];
                        _top3Donors[j] =  _top3Donors[j-1];
                        _top3Donors[j-1] = tAddress;
                    }
                }
                return;
            }
        }

        /*
         * 变量index，从0开始递增
         * 1、如果_top3Donors[index]是0地址，则将msg.sender赋给_top3Donors[index]
         * 2、不是0地址，判断msg.value是否不大于_top3Donors[index]的值，则index+1
         * 3、不是0地址，判断msg.value是否大于_top3Donors[index]的值，从index开始将_top3Donors[index]
         *    的值赋给_top3Donors[index+1]，最后将msg.sender赋给_top3Donors[index]
         */
        for (uint index = 0; index < 3;) {
            // 从左向右找第一个不为0地址，并且值小于msg.value
            if (_top3Donors[index] == address(0)) {
                _top3Donors[index] = msg.sender;
                break;
            } else {
                if (_donors[msg.sender] > _donors[_top3Donors[index]]) {
                    for(uint i =2; i >index; i--) {
                        _top3Donors[i] = _top3Donors[i-1];
                    }
                    _top3Donors[index] = msg.sender;
                    break;
                } else {
                    index++;
                }
            }
        }
    }

    function getTop3Donors() public view returns (address[3] memory) {
        return _top3Donors;
    }

    function getDonation(address from) public view returns (uint256) {
        return _donors[from];
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

    function getOwner() public view returns (address) {
        return _owner;
    }
}