/*
  - @Author: cc
  - @Date:   2023/12/21
  - @Description: 为了防止并发初始化，我们不能单单通过加锁的方式来实现Once
    使用doSlow方法应对并发初始化，并采用双检查机制(double-checking)
*/
package main

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	// double-checking
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
