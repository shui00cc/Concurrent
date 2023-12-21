/*
* @Author: cc
* @Date:   2023/12/19
* @Description: 通过Mutex实现线程安全的队列，避免入队出队时的data race
 */
package main

import (
	"fmt"
	"sync"
)

type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, v)
	fmt.Println("Enqueue:", v, "Now q:", q.data)
}

func (q *SliceQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		fmt.Println("no data can Dequeue")
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	fmt.Println("Dequeue:", v, "Now q:", q.data)
	return v
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	q := NewSliceQueue(5)
	go func() {
		defer wg.Done()
		q.Enqueue("cc")
	}()
	go func() {
		defer wg.Done()
		q.Enqueue(222)
	}()
	go func() {
		defer wg.Done()
		q.Dequeue()
	}()
	wg.Wait()
}
