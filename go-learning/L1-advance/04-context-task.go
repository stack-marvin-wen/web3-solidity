package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
 * @Description:构建一个支持取消和超时的任务系统
 * @Author: marvin.wensap@gmail.com
 * @Date: 2021-05-05 09:05:05
 * @LastEditTime: 2021-05-05 09:05:05
 */
type Task struct {
	ID       int
	Duration time.Duration
}

type TaskList struct {
	Tasks []Task
	mu    sync.Mutex
}

func (t1 *TaskList) AddTask(task Task) {
	t1.mu.Lock()
	defer t1.mu.Unlock()
	t1.Tasks = append(t1.Tasks, task)
}

func (t1 *TaskList) GetTasks() []Task {
	t1.mu.Lock()
	defer t1.mu.Unlock()
	return t1.Tasks
}
func (t1 *TaskList) RunTasks(ctx context.Context, task Task) error {
	ctx, cancel := context.WithTimeout(ctx, task.Duration)
	defer cancel()
	select {
	case <-time.After(task.Duration):
		fmt.Printf("任务 %d: 执行完成\n", task.ID)
		return nil
	case <-ctx.Done():
		fmt.Printf("任务 %d: 取消或超时: %v\n", task.ID, ctx.Err())
		return ctx.Err()
	}
}
func (t1 *TaskList) RunAllTasks(ctx context.Context) {
	tasks := t1.GetTasks()
	wg := sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			t1.RunTasks(ctx, t)
		}(task)
	}
	wg.Wait()
	fmt.Println("所有任务执行完成")
}
func main() {
	// 创建任务管理器
	tm := &TaskList{}

	// 添加任务
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 5; i++ {
		tm.AddTask(Task{
			ID:       i,
			Duration: time.Duration(rand.Intn(3)+1) * time.Second,
		})
	}

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动监听取消信号的goroutine
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n⚠️  手动触发取消信号")
		cancel()
	}()

	// 执行所有任务
	fmt.Println("开始执行任务...")
	tm.RunAllTasks(ctx)

	// 检查context状态
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("\n❌ 任务执行超时")
	} else if ctx.Err() == context.Canceled {
		fmt.Println("\n❌ 任务被手动取消")
	} else {
		fmt.Println("\n✅ 所有任务正常完成")
	}
}
