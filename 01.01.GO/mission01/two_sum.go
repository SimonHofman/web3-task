package main

import "fmt"

func twoSum(nums []int, target int) []int {
	fmt.Printf("nums: %v, target: %d, ", nums, target)
	index_arr := []int{}
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				index_arr = append(index_arr, i)
				index_arr = append(index_arr, j)
			}
		}
	}
	return index_arr
}

func main() {
	fmt.Printf("return: %v\n", twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Printf("return: %v\n", twoSum([]int{3, 2, 4}, 6))
	fmt.Printf("return: %v\n", twoSum([]int{3, 3}, 6))
}
