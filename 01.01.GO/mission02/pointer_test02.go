package main

import "fmt"

//题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
//考察点 ：指针运算、切片操作。

func doubleSlice(slice *[]int) {
	for i := 0; i < len(*slice); i++ {
		(*slice)[i] *= 2
	}
}

func main() {
	slice := []int{1, 2, 3, 4, 5}
	fmt.Println(slice)
	doubleSlice(&slice)
	fmt.Println(slice)
}
