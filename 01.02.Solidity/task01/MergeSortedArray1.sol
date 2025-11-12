// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract MergeSortedArray {
    function mergeSortedArray(uint[] memory arr1, uint[] memory arr2) public pure returns(uint[] memory) {
        uint newLength = arr1.length + arr2.length;
        uint[] memory arr3= new uint[](newLength);
        uint i = 0;
        uint j = 0;
        uint k = 0;
        while (i < arr1.length && j < arr2.length) {
            if (arr1[i]< arr2[j]) {
                arr3[k] = arr1[i];
                i++;
            } else {
                arr3[k] = arr2[j];
                j++;
            }
            k++;
        }
        while (i < arr1.length) {
            arr3[k] = arr1[i];
            i++;
            k++;
        }
        while (j < arr2.length) {
            arr3[k] = arr2[j];
            j++;
            k++;
        }
        return arr3;
    }
}
