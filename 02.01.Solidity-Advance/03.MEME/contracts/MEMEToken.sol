// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import {ERC20, ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

contract MEMEToken is ERC20Burnable {
    enum TransferMode {
        MODE_NORMAL,
        MODE_TRANSFER_RESTRICTED,
        MODE_TRANSFER_CONTROLED
    }

    TransferMode public  transferMode;
    address  public vestingContract;
    address  public pair;
    address public memeCore;

    error TransferRestricted();
    error TransferToTokenNotAllowed();
    error TransferNotAllowedToPair();
    error TransferNotAllowed();
    error onlyMemeCall();
    error ZeroAddress();

    event TransferModeChanged(TransferMode oldMode, TransferMode newMode);
    event VestingContractChanged(address vestingContract);
    event PairChanged(address pair);

    modifier onlyMeme() {
        if (msg.sender != memeCore) revert onlyMemeCall();
        _;
    }

    constructor(
        string memory name,
        string memory symbol,
        uint256 totalSuppply,
        address _meme
    ) ERC20(name, symbol) {
        if (_meme == address(0)) revert ZeroAddress();
        transferMode = TransferMode.MODE_TRANSFER_RESTRICTED;
        if (totalSuppply > 0) {
            _mint(_meme, totalSuppply);
        }
        memeCore = _meme;
    }

    function setTransferMode(TransferMode _mode) external onlyMeme {
        TransferMode oldMode = transferMode;
        transferMode = _mode;
        emit TransferModeChanged(oldMode, _mode);
    }

    function setVestingContract(address _vestingContract) external onlyMeme {
        if (_vestingContract == address(0)) revert ZeroAddress();
        vestingContract = _vestingContract;
        emit VestingContractChanged(_vestingContract);
    }

    function setPair(address _pair) external onlyMeme {
        if (_pair == address(0)) revert ZeroAddress();
        pair = _pair;
        emit PairChanged(_pair);
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {
        if (from == address(0) || to == address(0) || amount == 0) {
            return;
        }

        if (to == address(this)) {
            revert TransferToTokenNotAllowed();
        }

        if (from == vestingContract && vestingContract != address(0)) {
            return;
        }

        if (transferMode != TransferMode.MODE_NORMAL && to == pair && pair != address(0)) {
            revert TransferNotAllowedToPair();
        }

        if (transferMode == TransferMode.MODE_TRANSFER_RESTRICTED) {
            revert TransferRestricted();
        }
    }

    function _update(address  from , address  to, uint256 value) internal  override  {
        _beforeTokenTransfer(from, to, value);
        super._update(from, to, value);
    }
}
