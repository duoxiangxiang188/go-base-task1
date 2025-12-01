// SPDX-License-Identifier: MIT
pragma solidity ^0.8;


contract RomanToInt {
    function romanToInt(string memory roman) public view returns (uint256) {
        bytes memory  rb = bytes(roman);
        uint256 result = 0;
        uint256 length = rb.length;

        for (uint256 i = 0; i < length - 1; i++) {
            uint256 current = getValue(rb[i]);

            if (i < length - 1 && current < getValue(rb[i+1])){
                result -= current;

            } else {
                result += current;
            }

        }
        return result ;
    }

    function getValue(bytes1 c) internal pure  returns (uint256) {
        if(c == bytes1("I")) return  1;
        if(c == bytes1("V")) return  5;
        if(c == bytes1("X")) return  10;
        if(c == bytes1("L")) return  50;
        if(c == bytes1("C")) return  100;
        if(c == bytes1("D")) return  500;
        if(c == bytes1("M")) return  1000;

        revert("Invalid Roman character");
    }
}