package main

import "fmt"

func singleNumber(nums []int) int {

	for i := range nums {
		isUnique := true // 标记当前元素是否为唯一元素
		for j := 0; j < len(nums); j++ {
			if (i != j) && nums[i] == nums[j] {
				isUnique = false
				break // 找到重复元素，跳出内层循环
			}
		}
		// 如果是唯一元素，返回该元素
		if isUnique {
			return nums[i]
		}

	}
	return 0
}
func singleNumber1(nums []int) int {
	map1 := make(map[int]int)
	for i := range nums {
		map1[nums[i]]++
	}
	for key, value := range map1 {
		if value == 1 {
			//fmt.Println(key)
			return key
			break
		}
	}
	return 0
}

func main() {
	//1.简单定义输出测试
	nums1 := []int{2, 2, 1}
	fmt.Println(singleNumber(nums1))

	nums2 := []int{4, 1, 2, 1, 2}
	fmt.Println(singleNumber(nums2))

	nums3 := []int{1}
	fmt.Println(singleNumber(nums3))
	//2.二维数组遍历
	testCases := [][]int{
		{2, 2, 1},
		{4, 1, 2, 1, 2},
		{1},
	}
	for _, testCase := range testCases {
		result := singleNumber(testCase)
		fmt.Printf("数组 %v 中出现一次是元素是: %d\n", testCase, result)
	}
	for _, testCase := range testCases {
		result := singleNumber1(testCase)
		fmt.Printf("map方法数组 %v 中出现一次是元素是: %d\n", testCase, result)
	}

}
