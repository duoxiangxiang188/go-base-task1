package main

import (
	"fmt"
)

func isValid(s string) bool {
	//奇数长度直接返回false
	if len(s)%2 != 0 {
		return false
	}
	//定义括号映射表
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := make([]rune, 0, len(s)/2)

	for _, v := range s {

		switch v {
		case '(', '{', '[':
			//左括号入栈
			stack = append(stack, v)

		//右括号匹配栈顶元素
		case ')', '}', ']':
			//栈空时遇到有括号必无效
			if len(stack) == 0 {
				return false
			}
			//弹出栈顶元素以备之后匹配
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			//检查循环的右括号是否匹配
			if top != pairs[v] {
				return false
			}

		}

	}
	//最终栈必须为空才有效
	return len(stack) == 0
}

func main() {
	s1 := "()"

	fmt.Println(isValid(s1))
	s2 := "()[]{}"
	fmt.Println(isValid(s2))
	s3 := "(]"
	fmt.Println(isValid(s3))
	s4 := "([])"
	fmt.Println(isValid(s4))
	s5 := "([)]"
	fmt.Println(isValid(s5))
}
