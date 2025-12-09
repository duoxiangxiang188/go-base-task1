package main

import (
	"fmt"
)

func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		digits[i]++
		digits[i] %= 10
		if digits[i] != 0 {
			return digits
		}
	}
	return append([]int{1}, digits...)

}
func main() {
	strs1 := []int{1, 2, 3}

	fmt.Println(plusOne(strs1))
	strs2 := []int{4, 3, 2, 1}

	fmt.Println(plusOne(strs2))
	strs3 := []int{9}

	fmt.Println(plusOne(strs3))

}
