/*
* @Author: cc
* @Date:   2023/12/25
* @Description: 将channel当作流式管道(Stream)使用，并实现单进程的map-reduce方法
 */
package main

import "fmt"

// asStream 将提供的整数切片转换为一个流式管道，通过该管道可以进行流式处理。
// 使用 unbuffered channel，可通过关闭 done 通道来提前结束流式处理。
func asStream(done <-chan struct{}, values ...int) <-chan interface{} {
	s := make(chan interface{}) // 无缓冲的 channel
	go func() {
		defer close(s)
		for _, v := range values {
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}

// mapChan 对输入的流式管道进行映射操作，使用提供的映射函数 fn。
// 返回一个输出流式管道，其中包含映射后的数据。
func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) // 输出 channel
	if in == nil {                // 异常检查
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for v := range in { // 从输入管道读取数据，执行映射操作，并将结果发送到输出管道
			out <- fn(v)
		}
	}()
	return out
}

// reduce 对输入的流式管道进行缩减操作，使用提供的缩减函数 fn。
// 返回最终的缩减结果。
func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil { // 异常检查
		return nil
	}
	out := <-in         // 读取第一个元素作为初始值
	for v := range in { // 从输入管道读取数据，执行缩减操作，并更新结果
		out = fn(out, v)
	}
	return out
}

func main() {
	values := []int{1, 2, 3, 4, 5}
	in := asStream(nil, values...)
	// map操作：乘以10
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}
	// reduce操作：累加map的结果
	reduceFn := func(r, v interface{}) interface{} {
		return r.(int) + v.(int)
	}

	// 将输入流式管道进行映射和缩减操作，并输出最终结果
	sum := reduce(mapChan(in, mapFn), reduceFn)
	fmt.Println(sum)
}
