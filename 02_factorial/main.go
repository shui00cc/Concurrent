/*
  - @Author: cc
  - @Date:   2023/12/20
  - @Description:
    factorial()递归函数会不断获取读锁，并且全部子递归完成后才会解锁
    大约执行两次后写锁请求，此时发生死锁
    两个goroutine互相持有锁并等待，满足“writer依赖活跃的reader -> 活跃的reader依赖新来的reader -> 新来的reader依赖writer”的死锁条件
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.RWMutex

	// 模拟reader
	go func() {
		factorial(&mu, 10) // 计算10的阶乘
	}()

	// 模拟writer,稍微等待，然后制造一个调用Lock的场景
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Trying to acquire write lock")
		mu.Lock()
		fmt.Println("Write lock acquired")
		time.Sleep(100 * time.Millisecond)
		mu.Unlock()
		fmt.Println("Unlock")
	}()
	time.Sleep(1 * time.Second)
}

func factorial(m *sync.RWMutex, n int) int {
	if n < 1 { // 阶乘退出条件
		return 0
	}
	fmt.Println("Trying to acquire read lock")
	m.RLock()
	fmt.Println("Read lock acquired", n)
	defer func() {
		fmt.Println("RUnlock", n)
		m.RUnlock()
	}()
	time.Sleep(100 * time.Millisecond)
	return factorial(m, n-1) * n // 递归调用
}
