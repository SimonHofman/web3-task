// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";

import {Price} from "./libraries/RedBlackTreeLibrary.sol";
import {LibOrder, OrderKey} from "./libraries/LibOrder.sol";

abstract contract OrderValidator is Initializable, ContextUpgradeable, EIP712Upgradeable {
    bytes4 private constant EIP_1271_MAGIC_VALUE = 0x1626ba7e;

    uint256 private constant CANCELLED = type(uint256).max;

    mapping(OrderKey => uint256) public filledAmount;

    function __OrderValidator_init(string memory EIP712Name, string memory EIP712Version) internal onlyInitializing {
        __Context_init();
        __EIP712_init(EIP712Name, EIP712Version);
        __OrderValidator_init_unchained();
    }

    function __OrderValidator_init_unchained() internal onlyInitializing {}

    function _validateOrder(LibOrder.Order memory order, bool isSkipExpiry) internal view {
        require(order.maker != address(0), "OVa: miss maker");

        if (!isSkipExpiry) {
            require(order.expiry == 0 || order.expiry > block.timestamp, "OVa: expired");
        }

        require(order.salt != 0, "OVa: zero salt");

        if (order.side == LibOrder.Side.List) {
            require(order.nft.collection != address(0), "OVa: unsupported nft asset");
        } else if (order.side == LibOrder.Side.Bid) {
            require(Price.unwrap(order.price) > 0, "OVa: zero price");
        }
    }

    function _getFilledAmount(OrderKey orderKey) internal view returns (uint256 orderFilledAmount) {
        orderFilledAmount = filledAmount[orderKey];
        require(orderFilledAmount != CANCELLED, "OVa: order cancelled");
    }

    function _updateFilledAmount(uint256 newAmount, OrderKey orderKey) internal {
        require(newAmount != CANCELLED, "OVa: order cancelled");
        filledAmount[orderKey] = newAmount;
    }

    function _cancelOrder(OrderKey orderKey) internal {
        filledAmount[orderKey] = CANCELLED;
    }

    uint256[50] private __gap;
}
