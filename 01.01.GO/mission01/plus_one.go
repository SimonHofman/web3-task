package main

import "fmt"

func plusOne(digits []int) []int {
	var newDigit []int
	newDigit = make([]int, len(digits)+1)
	var tDigit int
	arr_len := len(digits)
	for i := len(newDigit) - 1; i >= 0; i-- {
		j := i - 1
		if i == arr_len {
			tDigit = digits[j] + 1
			newDigit[i] = tDigit % 10
		} else if i != 0 {
			tDigit = digits[j] + tDigit/10
			newDigit[i] = tDigit % 10
		} else if i == 0 && tDigit == 10 {
			newDigit[0] = tDigit / 10
		}
	}
	if newDigit[0] == 0 {
		return newDigit[1:]
	} else {
		return newDigit
	}
}

func main() {
	fmt.Println(plusOne([]int{1, 2, 3}))
	fmt.Println(plusOne([]int{4, 3, 2, 1}))
	fmt.Println(plusOne([]int{9}))
	fmt.Println(plusOne([]int{9, 9}))
}
