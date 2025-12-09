//SPDX-License-Identifier: MIT  
pragma solidity ^0.8.0;

contract MergeSortedArray {
    function merge(uint256[] calldata arr1, uint256[] calldata arr2)public pure  returns (uint256[] memory ){
        uint256 len1 = arr1.length;
        uint256 len2= arr2.length;
        uint256 num = len1+len2;
        uint256[] memory result = new uint256[](num);
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;

        while (i < len1 && j <len2) {
            if(arr1[i]<arr2[j]) {
                result[k++] = arr1[i++];
            } else {
                result[k++] = arr2[j++];
            }

        }

        while (i < len1) {
            result[k++] = arr1[i++];
        }
        while (j < len2) {
            result[k++] = arr2[j++];
        }
        return result;
    }
}