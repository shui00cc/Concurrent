/*
- @Author: cc
- @Date:   2023/12/27
- @Description: 官方的errgroup提供了.Group和.WithContext两种创建方式，提供了Wait和Go方法
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	// 执行三个子任务，第二个执行失败，通过cancel context取消监听的第一个进程
	g.Go(func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("exec #1")
		select {
		case <-ctx.Done():
			fmt.Println("#1 cancel, err =", ctx.Err())
		default:
		}
		return nil
	})
	g.Go(func() error {
		err := errors.New("#2 err")
		return err
	})
	g.Go(func() error {
		time.Sleep(3 * time.Second)
		fmt.Println("exec #3")
		return nil
	})
	// 等待任务完成
	if err := g.Wait(); err != nil {
		fmt.Println("failed:", err) // 返回子任务出现的第一个错误
	} else {
		fmt.Println("successfully exec all")
	}

}
