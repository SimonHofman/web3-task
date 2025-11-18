// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract BinarySearch1 {
    function binarySearch(uint256[] memory arr, uint256 x) public pure returns (int256) {
        int256 index = -1;
        uint256 right = arr.length -1;
        uint256 left = 0;
        while (left <= right) {
            uint256 mid =left + (right - left) / 2;
            if (arr[mid] == x) {
                index = int256(mid);
                break;
            }
            if (arr[mid] < x) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        return index;
    }
}
