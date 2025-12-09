package main

import (
	"fmt"
	"sync"
)

func receive(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("接收：%d\n", v)
	}
}
func send(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}

}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan int)
	go send(ch, &wg)
	go receive(ch, &wg)
	wg.Wait()
	fmt.Println("任务完成")

}
