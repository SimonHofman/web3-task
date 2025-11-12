// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract ReverseString {
    function reverseString(string memory input) public pure returns (string memory) {
        bytes memory inputBytes = bytes(input);
        uint inputBytesLength = inputBytes.length;
        bytes memory reversedBytes = new bytes(inputBytesLength);
        for (uint i = 0; i < inputBytesLength; i ++) {
            reversedBytes[inputBytesLength - i - 1] = inputBytes[i];
        }
        return string(reversedBytes);
    }
}
