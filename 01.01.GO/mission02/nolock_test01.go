package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var count int64 // 原子操作需要使用 int64 类型
	var wg sync.WaitGroup

	// 启动10个协称
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	// 等待所有携程结束
	fmt.Printf("count is %d\n", atomic.LoadInt64(&count))
}
