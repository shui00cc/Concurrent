/*
  - @Author: cc
  - @Date:   2023/12/22
  - @Description: 并发安全的map
    在读写锁的基础上，使用分片减少锁的粒度和持有时间
*/
package main

import (
	"sync"
)

var SHARED_COUNT = 32

type ConcurrentMap []*ConcurrentMapShared

type ConcurrentMapShared struct {
	items map[string]interface{}
	sync.RWMutex
}

func New() ConcurrentMap {
	m := make(ConcurrentMap, SHARED_COUNT)
	for i := 0; i < SHARED_COUNT; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

// GetShared 根据key哈希取模计算分片索引
func (m ConcurrentMap) GetShared(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARED_COUNT)]
}

// FNV hash
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	shared := m.GetShared(key)
	shared.Lock()
	shared.items[key] = value
	shared.Unlock()
}
func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	shared := m.GetShared(key)
	shared.Lock()
	value, existed := shared.items[key]
	shared.Unlock()
	return value, existed
}
