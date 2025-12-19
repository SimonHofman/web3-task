// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IVestingParams {
    enum VestingMode {
        BURN,
        ClIFF,
        LINEAR
    }

    struct VestingAllocation {
        uint256 amount;
        uint256 lauchTime;
        uint256 duration;
        VestingMode mode;
    }
}
