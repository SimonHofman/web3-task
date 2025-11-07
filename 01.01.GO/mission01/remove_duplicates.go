package main

import "fmt"

func removeDuplicates(nums *[]int) int {
	i, j := 0, 1
	// fmt.Printf("length: %d\n", len(*nums))
	for ; j < len(*nums); j++ {
		if (*nums)[i] != (*nums)[j] {
			(*nums)[i+1] = (*nums)[j]
			i++
		}
		// fmt.Printf("i: %d, j: %d, a: %v\n", i, j, *nums)
	}

	*nums = (*nums)[:i+1] // 截取前 i+1 个不重复元素

	return i + 1
}

//func removeDuplicates(nums *[]int) int {
//	if len(*nums) == 0 {
//		return 0
//	}
//
//	i := 0
//	for j := 1; j < len(*nums); j++ {
//		if (*nums)[i] != (*nums)[j] {
//			i++
//			(*nums)[i] = (*nums)[j]
//		}
//	}
//
//	*nums = (*nums)[:i+1] // 截取前 i+1 个不重复元素
//
//	return i + 1
//}

func main() {
	a := []int{1, 1, 2}
	fmt.Printf("origin a: %v, ", a)
	length := removeDuplicates(&a)
	fmt.Printf("new a: %v, len: %d\n", a, length)

	b := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Printf("origin a: %v, ", b)
	length = removeDuplicates(&b)
	fmt.Printf("new a: %v, len: %d\n", b, length)
}
