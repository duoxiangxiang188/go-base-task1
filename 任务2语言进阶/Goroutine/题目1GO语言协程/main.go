package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}) //控制交替的信号通道
	wg.Add(2)
	//打印奇数的协程
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数=", i)
			ch <- struct{}{} //通知偶数协程
			<-ch             //等待偶数协程通知
		}
	}()
	//打印偶数的协程
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			<-ch //等待奇数协程通知
			fmt.Println("偶数=", i)
			ch <- struct{}{} //通知奇数协程
		}
	}()
	wg.Wait()
	close(ch)
}
