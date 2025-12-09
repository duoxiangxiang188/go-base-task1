package main

import (
	"fmt"
)

func longestCommonPrefix(strs []string) string {
	//空字符串直接返回空
	if len(strs) == 0 {
		return ""
	}
	//取第i个字母进行比较
	for i := range strs[0] {
		//拿第一个字符串的第i个字母
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			//取到的第一个字符串的第i个字母长度高于要对比字符串长度，直接返回，或者字母不一致 返回之前的，如果第一个都不相等就返回空
			if i > len(strs[j]) || char != strs[j][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}
func main() {
	strs1 := []string{"fl1ower", "fl1ow", "fl1ight"}

	fmt.Println(longestCommonPrefix(strs1))
	strs2 := []string{"dog", "racecar", "car"}
	fmt.Println(longestCommonPrefix(strs2))

}
