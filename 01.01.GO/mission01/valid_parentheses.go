package main

import "fmt"

type byteStack []byte

func (s *byteStack) Push(value byte) {
	*s = append(*s, value)
}

func (s *byteStack) Pop() byte {
	if len(*s) == 0 {
		panic("stack is empty")
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value
}

func (s *byteStack) IsEmpty() bool {
	return len(*s) == 0
}

func validParentheses(s string) bool {
	fmt.Printf("%s: ", s)
	var stack byteStack
	char_arr := []byte(s)
	for _, item := range char_arr {
		switch item {
		case '(', '[', '{':
			stack.Push(item)
		case ')':
			if stack.IsEmpty() || stack.Pop() != '(' {
				return false
			}
		case ']':
			if stack.IsEmpty() || stack.Pop() != '[' {
				return false
			}
		case '}':
			if stack.IsEmpty() || stack.Pop() != '{' {
				return false
			}
		}
	}

	if stack.IsEmpty() {
		return true
	} else {
		return false
	}
}

func main() {
	fmt.Println(validParentheses("()"))
	fmt.Println(validParentheses("()[]{}"))
	fmt.Println(validParentheses("(]"))
	fmt.Println(validParentheses("([])"))
	fmt.Println(validParentheses("([)]"))
}
