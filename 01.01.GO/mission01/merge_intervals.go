package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func random_intervals() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		// 生成两个随机数作为区间的起点和终点
		start := rand.Intn(100)
		end := rand.Intn(100)

		// 确保起点小于等于终点
		if start > end {
			start, end = end, start
		}

		fmt.Printf("[%d, %d]\n", start, end)
	}
}

func merge_intervals(intervals [][]int) [][]int {
	// 创建副本避免修改原数组
	sorted_intervals := make([][]int, len(intervals))
	copy(sorted_intervals, intervals)

	// 按照区间的起始位置进行排序
	sort.Slice(sorted_intervals, func(i, j int) bool {
		return sorted_intervals[i][0] < sorted_intervals[j][0]
	})

	fmt.Printf("%v\n", sorted_intervals)

	for i := 0; i < len(sorted_intervals)-1; i++ {
		if sorted_intervals[i][1] >= sorted_intervals[i+1][0] {
			sorted_intervals[i][1] = max(sorted_intervals[i][1], sorted_intervals[i+1][1])
			sorted_intervals = append(sorted_intervals[:i+1], sorted_intervals[i+2:]...)
			i--
		}
	}

	return sorted_intervals
}

func main() {
	test := [][]int{{2, 6}, {1, 3}, {8, 10}, {15, 18}, {17, 20}}
	fmt.Printf("%v\n", merge_intervals(test))
}
