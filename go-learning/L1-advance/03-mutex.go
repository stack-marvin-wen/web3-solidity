package main

import (
	"fmt"
	"sync"
)

/*
 * @Description: 线程安全
 * @Author: marvin.wensap@gmail.com
 * @Date: 2021-05-05 09:05:05
 * @LastEditTime: 2021-05-05 09:05:05
 */

// 线程安全
// sync.Mutex是一个互斥锁，用于在多 goroutine 并发访问共享资源时实现互斥，保证同一时刻只有一个 goroutine 能访问被保护的代码区，防止数据竞争和并发安全问题。
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()         // 加锁，确保同一时刻只有一个 goroutine 能访问 count
	defer c.mu.Unlock() // 解锁，允许其他 goroutine 访问 count
	c.count++
}
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

type ReadWriteSafe struct {
	mu   sync.RWMutex
	data map[string]int
}

func (rw *ReadWriteSafe) Read(key string) int {
	rw.mu.RLock() // 读锁，允许多个 goroutine 同时读取
	defer rw.mu.RUnlock()
	return rw.data[key]
}
func (rw *ReadWriteSafe) Write(key string, value int) {
	rw.mu.Lock() // 写锁，确保同一时刻只有一个 goroutine 写入数据
	defer rw.mu.Unlock()
	rw.data[key] = value
}
func goroutineSync() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	sum := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			mu.Lock()
			sum += i
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}
func main() {
	goroutineSync()
}
