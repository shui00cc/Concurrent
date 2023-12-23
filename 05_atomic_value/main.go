/*
* @Author: cc
* @Date:   2023/12/23
* @Description: atomic.Value类型常用于配置变更场景，实现原子地读写对象
 */
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "chengdu",
		Addr:     "10.2.3.4",
		Count:    rand.Int31(),
	}
}

func main() {
	var cfg atomic.Value
	cfg.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})
	// 设置新config
	go func() {
		for {
			time.Sleep(time.Duration(1+rand.Int63n(3)) * time.Second)
			cfg.Store(loadNewConfig()) // Store写入新配置
			cond.Broadcast()           // 通知所有等待者配置已变更
		}
	}()
	// 读取新config
	go func() {
		cond.L.Lock()
		cond.Wait()              // 等待变更信号
		c := cfg.Load().(Config) // Load()读取新配置
		fmt.Printf("new config:%+v\n", c)
		cond.L.Unlock()
	}()
	time.Sleep(time.Second * 5)
}
