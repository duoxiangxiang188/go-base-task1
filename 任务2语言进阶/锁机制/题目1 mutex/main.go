package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 1; i <= 1000; i++ {
		c.value++
		// fmt.Println(c.value)
	}

}
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value

}

func main() {
	var counter Counter
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
			// fmt.Println(counter.value)
		}()
		// fmt.Println(counter.value)
	}

	wg.Wait()
	fmt.Println("最终计数器数值：", counter.value)

}
