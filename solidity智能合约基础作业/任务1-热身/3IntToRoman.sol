// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract IntToRoman {
    struct RomanPair {
        uint256 value;
        string symbol;
    }

    function intToRoman(uint256 num) public pure returns (string memory) {
        require(num > 0 && num <= 3999, "Number out of range(1-3999)");
        RomanPair[13] memory pairs = [
            RomanPair(1000, "M"),
            RomanPair(900, "CM"),
            RomanPair(500, "D"),
            RomanPair(400, "CD"),
            RomanPair(100, "C"),
            RomanPair(90, "XC"),
            RomanPair(50, "L"),
            RomanPair(40, "XL"),
            RomanPair(10, "X"),
            RomanPair(9, "IX"),
            RomanPair(5, "V"),
            RomanPair(4, "IV"),
            RomanPair(1, "I")
        ];
        bytes memory result = new bytes(0);

        for (uint256 i = 0; i < pairs.length; i++) {
            while (num >= pairs[i].value) {
                result = abi.encodePacked(result, pairs[i].symbol);
                num -= pairs[i].value;
            }
            if (num == 0) break;
        }

        return string(result);
    }
}
