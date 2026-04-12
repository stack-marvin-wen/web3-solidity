package main

import (
	"context"
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello")
}
func channelDemo() {
	fmt.Println("=== Channel示例 ===")
	// 无缓存channel
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- "Hello"
		ch <- "World"
	}()
	time.Sleep(3000 * time.Millisecond)
	for msg := range ch {
		fmt.Println("Received:", msg)
	}
	fmt.Println("=== Channel示例结束 ===")

}
func bufferedChannelDemo() {
	ch := make(chan int, 3) // make(chan int, 3) 的缓冲区只能容纳 3 个元素；你在同一个 goroutine里连续 ch <- 1、2、3、4，第四次写入时缓冲区已经满了，又没有并发的读取方法，所以写操作会一直阻塞。由于 main goroutine 被卡住，程序到不了后面的打印语句，最终触发运行时检测到 “all goroutines are asleep – deadlock!” 的 panic。要让第 4 次写入不阻塞，必须在写入时就有其他 goroutine 去消费
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4
	fmt.Println("Channel写入完成，读取...")
	fmt.Println("读取:", <-ch)
	fmt.Println("读取:", <-ch)
	fmt.Println("读取:", <-ch)
	fmt.Println("=== Buffered Channel示例结束 ===")
}
func bufferedChannelDemoNew() {
	fmt.Println("=== 生产消费模式下的Buffered Channel示例开始 ===")
	ch := make(chan int, 3)
	defer close(ch)

	go func() {
		for v := range ch {
			fmt.Println("channel读取(R1):", v)
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	go func() {
		for v := range ch {
			fmt.Println("channel读取(R2):", v)
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	for i := 0; i < 50; i++ {
		fmt.Println("channel写入:", i)
		ch <- i
	}
	fmt.Println("=== 生产消费模式下的Buffered Channel示例结束 ===")

}
func selectDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "ch1"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "ch2"
	}()
	// 随机选择一个就绪的channel
	select {
	case v := <-ch1:
		fmt.Println("ch1:", v)
	case v := <-ch2:
		fmt.Println("ch2:", v)
	}
}
func timeoutSelectDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "ch1"
	}()
	go func() {
		time.Sleep(3 * time.Second)
		ch2 <- "ch2"
	}()
	select {
	case v := <-ch1:
		fmt.Println("ch1:", v)
	case v := <-ch2:
		fmt.Println("ch2:", v)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	}
}

func loopCheckMuilpleChannel() {
	fmt.Println("监听开始")
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- fmt.Sprintf("ch1: %d", i)
			time.Sleep(time.Millisecond * 100)
		}
		close(ch1)
	}()
	go func() {
		for i := 0; i < 5; i++ {
			ch2 <- fmt.Sprintf("ch2: %d", i)
			time.Sleep(time.Millisecond * 150)
		}
		close(ch2)
	}()
	for ch1 != nil || ch2 != nil {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil
				continue
			}
			fmt.Println(v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println(v)
		default:
			fmt.Println("no data")
			time.Sleep(50 * time.Millisecond)

		}
	}
	fmt.Println("监听完成")
}
func quitSignChannelDemo() {
	jobs := make(chan int, 5)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case job := <-jobs:
				fmt.Println("处理job:", job)
			case <-quit:
				fmt.Println("收到退出信号，结束goroutine")
				return
			}
		}
	}()
	for i := 0; i < 10; i++ {
		jobs <- i
	}
	fmt.Println("发送任务完成")
	quit <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}
func closeChannelDemo() {
	ch := make(chan int)
	close(ch)
	select {
	case v, ok := <-ch:
		fmt.Printf("val: %d, ok: %v\n", v, ok)
	default:
		fmt.Println("没有数据")
	}
	fmt.Println()
}
func contextWithChannelDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- 42
	}()
	select {
	case v := <-ch:
		fmt.Println(v)
	case <-ctx.Done():
		fmt.Println("超时:", ctx.Err())
	}
}
func contextWithCancelSentDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("超时:", ctx.Err())
			cancel()
		default:
			fmt.Println("没有超时")
			time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second * 3)
	fmt.Println("信号发送")
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println()
}
func main() {
	// go sayHello() // goroutine调用
	// 方式1:等待执行结果，缺点无法固定每一个协程执行完成时间，sleep时间应该比协程执行时间长，所以sleep时间不好把控
	// time.Sleep(time.Second) // 主线程如果不等待将会看不到输出
	// 方式2:使用WaitGroup
	// var wg sync.WaitGroup
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1) // 添加一个任务
	// 	go func(i int) {
	// 		fmt.Println("Hello", i)
	// 		wg.Done() // 任务完成
	// 	}(i)
	// }
	// wg.Wait()

	// channelDemo()
	// bufferedChannelDemo()
	// bufferedChannelDemoNew()
	// selectDemo()
	// timeoutSelectDemo()
	// loopCheckMuilpleChannel()
	// quitSignChannelDemo()
	// closeChannelDemo()
	// contextWithChannelDemo()
	// contextWithCancelSentDemo()
}
