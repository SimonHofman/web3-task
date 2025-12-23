// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

interface IBondingCurveParams {
    struct BondingCurveParams {
        uint256 virtualBNBReserve;
        uint256 virtualTokenReserve;
        uint256 k;
        uint256 availableTokens;
        uint256 collectedBNB;
    }
}
