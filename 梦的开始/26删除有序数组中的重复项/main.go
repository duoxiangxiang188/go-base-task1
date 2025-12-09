package main

import (
	"fmt"
)

func removeDuplicates(nums []int) int {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] == nums[i+1] {
			nums = append(nums[:i], nums[i+1:]...)
			i--
		}
	}

	fmt.Println(len(nums), ", nums =", nums)
	return len(nums)
}
func main() {
	strs1 := []int{1, 1, 2}
	removeDuplicates(strs1)
	strs2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	removeDuplicates(strs2)

}
