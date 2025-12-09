package main

import (
	"fmt"
	"strconv"
)

func pnumber(i1 int) bool {
	if i1 < 0 { //负数直接排除
		return false
	}
	str := strconv.Itoa(i1)
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		if str[i] != str[j] {
			return false
		}

	}
	return true
}
func main() {
	var1 := 121121

	fmt.Println(pnumber(var1))
	x := -121
	fmt.Println(pnumber(x))
	x1 := 0
	fmt.Println(pnumber(x1))
}
