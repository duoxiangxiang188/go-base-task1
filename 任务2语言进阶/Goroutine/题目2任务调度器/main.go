package main

import (
	"fmt"
	"sync"
	"time"
)

// 定义任务类型：函数，无参数无返回值
type Task func()

type NamedTask struct {
	Name string
	Fn   Task
}

// 调度器接收NameTask切片
func RunTasks(tasks []NamedTask) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, nt := range tasks {
		go func(name string, t Task) {
			defer wg.Done()
			start := time.Now()
			t()
			eclapsed := time.Since(start)
			fmt.Printf("任务【%s】完成，用时: %v \n", name, eclapsed)
		}(nt.Name, nt.Fn)
	}
	wg.Wait()
}

func main() {
	tasks := []NamedTask{
		{Name: "任务1", Fn: func() { time.Sleep(1 * time.Second) }},
		{Name: "任务2", Fn: func() { time.Sleep(2 * time.Second) }},
	}
	RunTasks(tasks)
}
