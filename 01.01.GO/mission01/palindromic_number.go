package main

import "fmt"

func palindromicNumber(num int) bool {
	return num == reverseNumber(num)
}

func reverseNumber(num int) int {
	var reverseNum int
	for num > 0 {
		reverseNum = reverseNum*10 + num%10
		num = num / 10
	}
	return reverseNum
}

func main() {
	fmt.Println(palindromicNumber(121))
	fmt.Println(palindromicNumber(1))
	fmt.Println(palindromicNumber(12))
}
