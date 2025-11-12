// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract RomanToInt1 {
    function fomanToInt(string memory roman) public pure returns (uint) {
        uint result = 0;
        bytes memory romanBytes = bytes(roman);
        uint length = romanBytes.length;
        for (uint i = 0; i < length;) {
            if (romanBytes[i] == "M") {
                result += 1000;
            } else if (romanBytes[i] == "D") {
                result += 500;
            } else if (romanBytes[i] == "C" && i + 1 < length) {
                if (romanBytes[i + 1] == "M") {
                    result += 900;
                    i++;
                } else if (romanBytes[i + 1] == "D") {
                    result += 400;
                    i++;
                } else {
                    result += 100;
                }
            } else if (romanBytes[i] == "L") {
                result += 50;
            } else if (romanBytes[i] == "X" && i + 1 < length) {
                if (romanBytes[i + 1] == "C") {
                    result += 90;
                    i++;
                } else if (romanBytes[i + 1] == "L") {
                    result += 40;
                    i++;
                } else {
                    result += 10;
                }
            } else if (romanBytes[i] == "V") {
                result += 5;
            } else if (romanBytes[i] == "I" && i + 1 < length) {
                if (romanBytes[i + 1] == "X") {
                    result += 9;
                    i++;
                } else if (romanBytes[i + 1] == "V") {
                    result += 4;
                    i++;
                } else {
                    result += 1;
                }
            } else if (romanBytes[i] == "I"){
                result += 1;
            }
            i++;
        }
        return (result);
    }
}
