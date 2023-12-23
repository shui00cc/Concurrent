/*
  - @Author: cc
  - @Date:   2023/12/23
  - @Description: channel常用于任务编排、消息传递与通知、Select等
    使用channel实现任务编排：
    有四个goroutine，编号为 1、2、3、4，每秒钟会有一个 goroutine 打印出它自己的编号，输出的编号总是按照 1、2、3、4、1、2、3、4…的顺序打印出来
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

const taskCount = 4

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	// 启动goroutine, num为编号, 从0开始
	for num := 0; num < taskCount; num++ {
		wg.Add(1)
		go printTask(num, ch, &wg)
	}
	// 通过channel发送信号，控制goroutine的执行顺序
	for i := 0; i <= 10; i++ {
		ch <- i % (taskCount)
		time.Sleep(time.Second) // 每秒发送一个信号
	}

	close(ch) // 关闭channel，通知所有goroutine退出
	wg.Wait()
}
func printTask(num int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case n, ok := <-ch:
			if !ok {
				return // channel关闭时退出
			}
			if n == num { // 当收到自己的编号时打印
				fmt.Println(num + 1)
			}
		}
	}
}
