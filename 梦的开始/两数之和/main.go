package main

import (
	"fmt"
)

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		a := target - nums[i]
		for j := 0; j < len(nums); j++ {
			if i != j && a == nums[j] {
				return []int{i, j}
			}
		}

	}
	return []int{}
}

// func twoSum(nums []int, target int) []int {
// 	numMap := make(map[int]int)
// 	for i, num := range nums {
// 		n := target - num
// 		if idx, exists := numMap[n]; exists {
// 			return []int{idx, i}
// 		}
// 		numMap[num] = i
// 	}
// 	return nil
// }

func main() {
	// 测试用例
	intervals1 := []int{2, 7, 11, 15}
	target1 := 9

	fmt.Println(twoSum(intervals1, target1))
	intervals2 := []int{3, 2, 4}
	target2 := 6

	fmt.Println(twoSum(intervals2, target2))
	intervals3 := []int{3, 3}
	target3 := 6

	fmt.Println(twoSum(intervals3, target3))

}
