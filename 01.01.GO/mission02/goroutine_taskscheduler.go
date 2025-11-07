package main

import (
	"fmt"
	"sync"
	"time"
)

// task 任务结构体
type Task struct {
	Name string
	Func func() error
}

// TaskResult 任务执行结果
type TaskResult struct {
	Name     string
	Duration time.Duration
	Error    error
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	tasks   []Task
	results chan TaskResult
	wg      sync.WaitGroup
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:   make([]Task, 0),
		results: make(chan TaskResult, 10),
	}
}

// AddTask 添加任务
func (ts *TaskScheduler) AddTask(name string, taskFunc func() error) {
	ts.tasks = append(ts.tasks, Task{Name: name, Func: taskFunc})
}

// Execute 并发执行所有任务
func (ts *TaskScheduler) Execute() []TaskResult {
	// 启动所有任务
	for _, task := range ts.tasks {
		ts.wg.Add(1)
		go ts.runTask(task)
	}

	// 启动结果收集携程
	go func() {
		ts.wg.Wait()
		close(ts.results)
	}()

	// 收集结果
	var results []TaskResult
	for result := range ts.results {
		results = append(results, result)
	}

	return results
}

// runTask 执行单个任务
func (ts *TaskScheduler) runTask(task Task) {
	defer ts.wg.Done()

	start := time.Now()
	err := task.Func()
	duration := time.Since(start)

	// 发送结果到通道
	ts.results <- TaskResult{
		Name:     task.Name,
		Duration: duration,
		Error:    err,
	}
}

// 示例任务函数
func task1() error {
	time.Sleep(2 * time.Second)
	fmt.Println("Task 1 completed")
	return nil
}

func task2() error {
	time.Sleep(1 * time.Second)
	fmt.Println("Task 2 completed")
	return nil
}

func task3() error {
	time.Sleep(3 * time.Second)
	fmt.Println("Task 3 completed")
	return nil
}

func main() {
	// 创建调度器
	scheduler := NewTaskScheduler()

	// 添加任务
	scheduler.AddTask("Download File", task1)
	scheduler.AddTask("Process Data", task2)
	scheduler.AddTask("Upload Result", task3)

	// 执行任务并获取结果
	fmt.Println("Starting task execution...")
	results := scheduler.Execute()

	// 打印执行结果
	fmt.Println("\nTask Execution Results:")
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("Task: %s, Duration: %v, Error: %v\n", result.Name, result.Duration, result.Error)
		} else {
			fmt.Printf("Task: %s, Duration: %v, Status: Success\n", result.Name, result.Duration)
		}
	}
}
