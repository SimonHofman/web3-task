// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract IntToRoman2 {
    uint[] numbers = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
    string[] symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];

    function intToRoman(uint num) public view returns (string memory) {
        string memory roman = "";

        for (uint i = 0; i < numbers.length; ) {
            if (num >= numbers[i]) {
                num -= numbers[i];
                roman = string.concat(roman, symbols[i]);
            } else if (num > 0){
                i++;
            } else {
                break;
            }
        }

        return roman;
    }
}

