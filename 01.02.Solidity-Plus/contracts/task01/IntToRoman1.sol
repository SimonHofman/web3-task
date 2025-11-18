// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract IntToRoman1 {
    function intToRoman(uint num) public pure returns (string memory) {
        string memory roman = "";

        for (; num > 0; ) {
            if (num >= 1000) {
                num -= 1000;
                roman = string.concat(roman, "M");
            } else if (num >= 900) {
                num -= 900;
                roman = string.concat(roman,"CM");
            } else if (num >= 500) {
                num -= 500;
                roman = string.concat(roman,"D");
            } else if (num >= 400) {
                num -= 400;
                roman = string.concat(roman,"CD");
            } else if (num >= 100) {
                num -= 100;
                roman = string.concat(roman,"C");
            } else if (num >= 90) {
                num -= 90;
                roman = string.concat(roman,"XC");
            } else if (num >= 50) {
                num -= 50;
                roman = string.concat(roman,"L");
            } else if (num >= 40) {
                num -= 40;
                roman = string.concat(roman,"XL");
            } else if (num >= 10) {
                num -= 10;
                roman = string.concat(roman,"X");
            } else if (num >=9) {
                num -= 9;
                roman = string.concat(roman,"IX");
            } else if (num >= 5) {
                num -= 5;
                roman = string.concat(roman,"V");
            } else if (num >= 4) {
                num -= 4;
                roman = string.concat(roman,"IV");
            } else if (num >= 1) {
                num -= 1;
                roman = string.concat(roman,"I");
            }
        }
        return roman;
    }
}

