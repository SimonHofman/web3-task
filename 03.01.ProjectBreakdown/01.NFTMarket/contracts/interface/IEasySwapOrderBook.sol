// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {OrderKey, Price, LibOrder} from "../libraries/LibOrder.sol";

interface IEasySwapOrderBook {
    function makeOrders(LibOrder.Order[] calldata newOrders) external payable returns (OrderKey[] memory newOrderKeys);

    function cancelOrders(OrderKey[] calldata orderKeys) external returns (bool[] memory successes);

    function editOrders(
        LibOrder.EditDetail[] calldata editDetails
    ) external payable returns (OrderKey[] memory newOrderKeys);

    function matchOrder(
        LibOrder.Order calldata sellOrder,
        LibOrder.Order calldata buyOrder
    ) external payable;

    function matchOrders(
        LibOrder.MatchDetail[] calldata matchDetails
    ) external payable returns (bool[] memory successes);
}
