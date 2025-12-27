// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import {LibPayInfo} from "./libraries/LibPayInfo.sol";

abstract  contract ProtocolManager is Initializable, OwnableUpgradeable {
    uint128 public protocalShare;

    event LogUpdatedProtocalShare(uint128 indexed newProtocolShare);

    function __ProtocolManager_init(uint128 newProtocolShare) internal  onlyInitializing {
        __ProtocolManager_init_unchained(newProtocolShare);
    }

    function __ProtocolManager_init_unchained(uint128 newProtocolShare) internal  onlyInitializing {
        _setProtocolShare(newProtocolShare);
    }

    function setProtocolShare(uint128 newProtocolShare) external  onlyOwner {
        _setProtocolShare(newProtocolShare);
    }

    function _setProtocolShare(uint128 newProtocalShare) internal  {
        require(
            newProtocalShare <= LibPayInfo.MAX_PROTOCOL_SHARE,
            "PM: exceed max protocol share"
        );
        protocalShare = newProtocalShare;
        emit LogUpdatedProtocalShare(newProtocalShare);
    }

    // 存储插槽预留，用于合约升级时保持存储布局兼容性：
    uint256[50] private  __gap;
}
