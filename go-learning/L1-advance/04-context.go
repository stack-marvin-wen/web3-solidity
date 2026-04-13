package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 1. 取消控制示例
// 2. 超时控制示例
// 3. 截止时间示例
func cancelableDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine收到取消信号:", ctx.Err())
				return
			default:
				fmt.Println("没有取消")
				time.Sleep(time.Second * 1)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	fmt.Println("发送取消信号(等待两秒)")
	cancel()
	time.Sleep(time.Second * 2)
}
func timeoutDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine收到取消信号:", ctx.Err())
				return
			default:
				fmt.Println("没有取消")
				time.Sleep(time.Second * 1)
			}
		}
	}()
	time.Sleep(time.Second * 5)
	fmt.Println("主函数结束")
	time.Sleep(time.Second * 2)
	fmt.Println()
}
func deadlineContextDemo() {
	fmt.Println("=== 截止时间Context ===")

	// 设置3秒后的截止时间
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 检查剩余时间
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("截止时间: %v, 剩余: %v\n", d, time.Until(d))
	}

	// 等待超过截止时间
	time.Sleep(4 * time.Second)

	select {
	case <-ctx.Done():
		fmt.Println("已超过截止时间:", ctx.Err())
	default:
		fmt.Println("未超时")
	}
}
func valueContextDemo() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "name", "zhangsan")
	ctx = context.WithValue(ctx, "age", 50)
	precessValueRequest(ctx)
}
func precessValueRequest(ctx context.Context) {
	if name := ctx.Value("name"); name != nil {
		fmt.Println("name:", name)
	}
	if age := ctx.Value("age"); age != nil {
		fmt.Println("age:", age)
	}
}

func cascadeCancelDemo() {
	parentCtx, parentCanel := context.WithCancel(context.Background())
	defer parentCanel()
	childCtx1, childCanel1 := context.WithCancel(parentCtx)
	defer childCanel1()
	childCtx2, childCanel2 := context.WithCancel(parentCtx)
	defer childCanel2()

	go worker(childCtx1, "child1")
	go worker(childCtx2, "child2")

	time.Sleep(3 * time.Second)
	fmt.Println("父Context发送取消信号")
	parentCanel()
	time.Sleep(2 * time.Second)

}
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("取消操作:", name, ":", ctx.Err())
			return
		default:
			fmt.Println("正在处理:", name)
			time.Sleep(time.Second)
		}
	}
}
func httpRequestDemo() {
	ctx, cannel := context.WithTimeout(context.Background(), time.Second*2)
	defer cannel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("请求超时:", err)
			return
		} else {
			fmt.Println("请求失败:", err)
		}
	}
	resp.Body.Close()
	fmt.Println("请求成功")
}

func multiWorkerDemo() {
	ctx, cannel := context.WithCancel(context.Background())
	defer cannel()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Worker收到取消信号:", ctx.Err())
					return
				default:
					fmt.Println("Worker正在工作")
					time.Sleep(time.Millisecond * 500)
				}
			}
		}()
	}
	time.Sleep(time.Second * 3)
	fmt.Println("取消所有Worker")
	cannel()
	wg.Wait()
	fmt.Println("所有Worker已取消")
}
func main() {
	// cancelableDemo()
	// timeoutDemo()
	// deadlineContextDemo()
	// valueContextDemo()
	// cascadeCancelDemo()
	// httpRequestDemo()
	multiWorkerDemo()
}
