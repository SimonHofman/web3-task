// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {IMEMEFactory} from "./interfaces/IMEMEFactory.sol";
import {MEMEToken} from "./MEMEToken.sol";

contract MEMEFactory is IMEMEFactory, AccessControl {
    bytes32 public constant DEPLOYER_ROLE = keccak256("DEPLOYER_ROLE");
    address public MEME;
    constructor(address _admin) {
        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
    }

    function deployToken(
        string memory name,
        string memory symbol,
        uint256 totalSupply,
        uint256 timestamp,
        uint256 nonce
    ) external onlyRole(DEPLOYER_ROLE) returns (address) {
        bytes32 salt = keccak256(abi.encodePacked(name, symbol, totalSupply, MEME, timestamp, nonce));
        MEMEToken token = new MEMEToken{salt: salt}(name, symbol, totalSupply, MEME);
        emit TokenDeployed(address(token), name, symbol, totalSupply, msg.sender);
        return address(token);
    }

    function predictTokenAddress(
        string memory name,
        string memory symbol,
        uint256 totalSupply,
        address owner,
        uint256 timestamp,
        uint256 nonce
    ) external view returns (address) {
        bytes32 salt = keccak256(abi.encodePacked(name, symbol, totalSupply, MEME, timestamp, nonce));

        bytes32 hash = keccak256(abi.encodePacked(
            bytes1(0xff),
            address(this),
            salt,
            keccak256(abi.encodePacked(
                type(MEMEToken).creationCode,
                abi.encode(name, symbol, totalSupply, owner)
            ))
        ));

        return address(uint160(uint256(hash)));
    }

    function setMEME(address  _MEME) external  onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_MEME != address (0), "ZeroAddress");
        _grantRole(DEPLOYER_ROLE, _MEME);
        MEME = _MEME;
    }
}
