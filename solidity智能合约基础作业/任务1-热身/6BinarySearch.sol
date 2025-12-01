// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract BinarySearch {
    function binarySearch(
        uint256[] calldata num,
        uint256 target
    ) public pure returns (int256) {
        if (num.length == 0) {
            return -1;
        }
        uint256 left = 0;
        uint256 right = num.length - 1;

        while (left <= right) {
            uint256 mid = left + (right - left) / 2;
            if (num[mid] == target) {
                return int256(mid);
            } else if (num[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        return -1;
    }
}
