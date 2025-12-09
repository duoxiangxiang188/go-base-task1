package main

import "fmt"

//指针加10
func pointPlus(p *int) {
	*p += 10
}

//指针切片乘以2
func pointMul(ps *[]int) {
	// for i, _ := range *ps {
	// 	(*ps)[i] *= 2
	// }
	for i, v := range *ps {
		(*ps)[i] = v * 2
	}
}
func main() {
	//
	v := 1
	pointPlus(&v)
	fmt.Println(v)

	v1 := []int{1, 3, 2, 5, 4}
	pointMul(&v1)
	fmt.Println(v1)
}
