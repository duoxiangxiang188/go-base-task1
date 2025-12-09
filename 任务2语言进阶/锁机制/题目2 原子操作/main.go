package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value int32
}

func (c *Counter) Inc() {

	for i := 1; i <= 1000; i++ {
		atomic.AddInt32(&c.value, 1)
		// fmt.Println(c.value)
	}

}
func (c *Counter) Value() int32 {
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
