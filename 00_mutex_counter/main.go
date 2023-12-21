/*
* @Author: cc
* @Date:   2023/12/19
* @Description: 使用Mutex.Lock()实现计数器
 */
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	sync.Mutex
	count uint64
}

// Incr 方法： 计数加一，内部使用互斥锁保护
func (c *Counter) Incr() {
	c.Lock()
	c.count++
	c.Unlock()
}

func main() {
	// 封装好的计数器
	var counter Counter

	var wg sync.WaitGroup
	wg.Add(10)

	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// 每个goroutine需要加100次
			for j := 0; j < 100; j++ {
				counter.Incr()
			}
		}()
	}
	// 等待完成
	wg.Wait()
	fmt.Println("count:", counter.count)
}
