package main

import "fmt"

func singleNumber(nums []int) int {
	count := make(map[int]int)
	for _, num := range nums {
		count[num]++
	}
	for k, v := range count {
		if v == 1 {
			return k
		}
	}
	return count[0]
}

func main() {
	fmt.Println(singleNumber([]int{2, 2, 1}))
	fmt.Println(singleNumber([]int{4, 1, 2, 1, 2}))
	fmt.Println(singleNumber([]int{1}))
}
