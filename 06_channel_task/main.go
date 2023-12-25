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
	"time"
)

var taskCount = 5

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch       // 取令牌
		fmt.Println(id + 1) // 从1开始打印
		time.Sleep(time.Second)
		nextCh <- token
	}
}

func main() {
	// 创建 taskCount 个channel
	var chs []chan Token
	for i := 0; i < taskCount; i++ {
		chs = append(chs, make(chan Token))
	}
	// 启动newWorker
	for id := 0; id < taskCount; id++ {
		go newWorker(id, chs[id], chs[(id+1)%taskCount])
	}
	// 初始token给第一个chan
	chs[0] <- Token{}
	select {}
}
