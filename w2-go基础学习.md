# Week 2: go语言基础学习
## 1. go数据类型
- 整型：int8, int16, int32, int64, int
- 浮点型：float32, float64
- 布尔型：bool
- 字符串：string
- 数组：固定长度的元素集合，类型必须相同。
- 切片：动态长度的元素集合，类型必须相同。
- map：键值对集合，键必须唯一，值可以重复。
## 2. 控制语句
- if 语句：if 条件 { 语句 }
- switch 语句：switch 值 { case 值: 语句; break; default: 语句; }
- for 语句：for 循环条件 { 循环体 }
- 循环控制语句：break, continue
- defer 语句：defer 函数名()
## 3. 函数：
```
func 函数名(参数列表) 返回值列表 { 
    函数体 
}
```
### 3.1 闭包

闭包（Closure）是指一个函数“捕获”了其外部作用域的变量，并可以在函数体内使用这些变量。Go 语言中的闭包常用于函数作为返回值或参数时，能够记住并操作其创建时的上下文变量。

**示例：**
```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

f := adder()
fmt.Println(f(10)) // 输出 10
fmt.Println(f(20)) // 输出 30
fmt.Println(f(30)) // 输出 60
```
上例中，返回的匿名函数引用了外部变量 sum，每次调用都会累加并返回结果，这就是闭包的效果。

**执行过程与内存示意图：**

1. 调用 `adder()` 时，创建了变量 `sum`，并返回一个匿名函数，这个匿名函数持有对 `sum` 的引用。
2. 变量 `f` 保存了这个匿名函数，每次调用 `f(x)`，其实就是调用这个匿名函数，操作的都是同一个 `sum` 变量。
3. 即使 `adder()` 已经返回，`sum` 依然不会被销毁，因为闭包引用了它。

```
内存示意图：

+-------------------+         +----------------------+
|       f           | ----->  |  匿名函数(闭包)      |
+-------------------+         +----------------------+
                                   |
                                   v
                             +-----------+
                             |   sum=30  |  // sum在多次调用后累加
                             +-----------+
```

每次调用 `f(x)`，闭包内部的 `sum` 都会被更新和保留，实现了“记忆”效果。

**使用场景与举例：**

1. **数据封装与状态保持**  
   闭包可以隐藏内部状态，只暴露操作接口。例如累加器、计数器等。

   ```go
   func counter() func() int {
       count := 0
       return func() int {
           count++
           return count
       }
   }

   c := counter()
   fmt.Println(c()) // 1
   fmt.Println(c()) // 2
   ```

2. **回调函数**  
   在异步编程或事件处理中，闭包常用于保存上下文信息。

   ```go
   func process(data []int, handler func(int) int) []int {
       result := []int{}
       for _, v := range data {
           result = append(result, handler(v))
       }
       return result
   }

   factor := 2
   doubled := process([]int{1, 2, 3}, func(x int) int {
       return x * factor // handler 闭包捕获了 factor
   })
   fmt.Println(doubled) // [2 4 6]
   ```

3. **工厂函数**  
   用于生成带有自定义行为的函数。

   ```go
   func makeMultiplier(factor int) func(int) int {
       return func(x int) int {
           return x * factor
       }
   }

   double := makeMultiplier(2)
   triple := makeMultiplier(3)
   fmt.Println(double(5)) // 10
   fmt.Println(triple(5)) // 15
   ```

闭包让函数拥有“记忆”能力，适合需要保存状态或上下文的场景。
### 3.2 柯里化

柯里化（Currying）是一种将带有多个参数的函数，转换为一系列只带一个参数的函数的技术。每次调用返回一个新函数，直到所有参数都被传递，最终执行计算。柯里化常用于函数式编程，便于函数复用和组合。

**Go 语言中的柯里化示例：**
```go
func add(a int) func(int) int {
    return func(b int) int {
        return a + b
    }
}

add5 := add(5)      // 返回一个新函数，参数为b
fmt.Println(add5(3)) // 输出 8
fmt.Println(add(2)(4)) // 输出 6
```
上例中，`add` 是一个柯里化函数，先传入参数 `a`，返回一个新函数，再传入参数 `b`，最终返回 `a + b` 的结果。

**柯里化的优点：**

1. **函数复用**：柯里化使得函数的部分应用变得简单，便于创建高阶函数。
2. **延迟执行**：参数可以分多次传递，支持更灵活的调用方式。
3. **提高可读性**：通过明确的参数分离，提升代码的可读性和可维护性。

**柯里化的应用场景：**

- **函数组合**：将多个处理步骤的函数，柯里化后可以方便的组合成一个新函数。
- **事件处理**：在事件驱动的编程中，柯里化可以简化事件处理函数的参数传递。
- **配置设置**：对于需要多步配置的函数，柯里化可以让每一步的意图更加明确。

柯里化是函数式编程中的一个重要概念，通过将函数参数分解为多个单一参数的函数调用，提供了更高的灵活性和可组合性。

## 4. 接口(interface{})
在Go语言中，接口（interface）用来描述一组方法签名，任何类型只要实现了接口中声明的全部方法，就被视为满足该接口。Go不要求显式声明“implements”，这被称为隐式实现：类型的结构和功能自然地决定了它属于哪个接口。这样的设计让代码更松耦合，也便于我们在不改动已有类型的前提下，为它们赋予新的行为。
### 4.1 接口定义
```go
// 定义接口
type Shape interface {
    Area() float64
    Perimeter() float64
}
```
接口定义了一个类型，该类型必须实现接口中的所有方法。
### 4.2 接口实现
1. **结构体实现**
```go
type Rectangle struct {
    width  float64
    height float64
}
func (r Rectangle) Area() float64 {
    return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.width + r.height)
}
```
2. **方法实现**
```go
type Circle struct {
    radius float64
}
func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}
func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.radius
}
func (c *Circle) SetRadius(radius float64) {
    c.radius = radius
}
```
3. **接口调用**
```go
func main() {
    var s Shape
    s = Rectangle{width: 3, height: 4}
    fmt.Printf("Rectangle Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())

    s = Circle{radius: 5}
    fmt.Printf("Circle Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}
```
代码解释: 接口Shape定义了两个方法Area()和Perimeter()，接口的实现者必须实现这两个方法。 

接口Shape的实现者可以是结构体Rectangle和结构体Circle，它们分别实现了接口Shape中的方法Area()和Perimeter()。

接口的调用通过接口变量s进行，s可以指向任何实现了Shape接口的类型实例。在main函数中，s先指向一个Rectangle实例，然后指向一个Circle实例，通过调用Area()和Perimeter()方法来计算并输出结果。

### 4.3 空类型与类型断言

```go
var i interface{}
```
这里将i定义为interface{}，即空接口。空接口可以存储任何类型的值，因此空接口可以存储任何类型的值。
- 空接口使用场景
    - fmt.Println、log.Print 等变长参数需要接受任意类型。

    - encoding/json 反序列化到 map[string]interface{}，让我们按键提取动态字段。

    - 抽象事件总线、消息队列时，可以先用空接口承载载荷，后续配合断言或注册回调处理。
- 最佳实践提醒
    - 提前约定空接口中可能出现的具体类型，配合文档或命名约束。
    - 优先考虑具体类型或泛型（Go 1.18+），空接口是兜底方案。
    - 使用类型断言或 type switch 时，一定要处理失败分支，避免 panic。
空接口（interface{}）是一种万能类型，它可以存储任何类型的值。但是，空接口的用法和具体类型有差异。
1. **类型断言**  
类型断言（Type Assertion）用于将一个接口变量转换为具体类型。
语法：
```go
var i interface{}
if v, ok := i.(T); ok {}
v := i.(T)
```
第二行代码的作用是判断i是否为T类型，如果为T类型，则返回T类型的值，并赋值给变量v，否则返回false。

这里的T是具体类型，比如int、string等。第三行代码的作用是将i转换为具体类型T，并赋值给变量v。
```go
func emptyInterfaceDemo() {
	var i interface{}
	i = 42
	// 类型断言
	v, ok := i.(int) // 类型断言，检查i是否是int类型
	if ok {
		fmt.Printf("i的值是: %d\n", v)
	} else {
		fmt.Println("i不是int类型")
	}
	// 在switch中使用i.(type)判断i的具体类型
	switch v := i.(type) {
	case int:
		fmt.Printf("i是int类型，值为: %d\n", v)
	case string:
		fmt.Printf("i是string类型，值为: %s\n", v)
	default:
		fmt.Printf("i是其他类型，值为: %v\n", v)
	}
	
}
```
### 4.4 多个interface{}组合
```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
type Writer interface {
	Write(p []byte) (n int, err error)
}
type ReadWriter interface {
	Reader
	Writer
}
type File struct {
	name string
}

func (f *File) Reader(data []byte) (n int, err error) {
	return 0, nil
}
func (f *File) Write(data []byte) (n int, err error) {
	return 0, nil
}
```
## 5. goroutine/channel/select/mutex/sync/
### 5.1 goroutine
- goroutine 是 go 语言中的并发原语，它允许多个函数在相同的地址空间内并发执行。
- goroutine内部运行时，会创建一个栈，栈的容量大小由系统决定。
- goroutine 的创建和销毁由 go 运行时管理，不需要用户管理。
- goroutine也称为协线程, 协程
- 创建goroutine, 使用go关键字, 跟普通函数调用一样
- 协程调用内存示意图：

```
主线程（main goroutine）
    |
    +---> go func(i int) {...} // 子协程1
    |
    +---> go func(i int) {...} // 子协程2
    |
    +---> ...                  // 更多子协程

每个 goroutine 都有独立的栈空间，和主线程共享堆内存。
主线程通过 sync.WaitGroup 等待所有子协程完成。

+-------------------+      +-------------------+
|   主线程栈(main)  |      |   子协程栈(goroutine) |
+-------------------+      +-------------------+
           |                        |
           |                        |
           +--------共享堆内存-------+
```

- 示例代码：
```go
func sayHello() {
	fmt.Println("Hello")
}
func main() {
	// go sayHello() // goroutine调用
	// 方式1:等待执行结果，缺点无法固定每一个协程执行完成时间，sleep时间应该比协程执行时间长，所以sleep时间不好把控
	// time.Sleep(time.Second) // 主线程如果不等待将会看不到输出
	// 方式2:使用WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1) // 添加一个任务
		go func(i int) {
			fmt.Println("Hello", i)
			wg.Done() // 任务完成
		}(i)
	}
	wg.Wait()
}
```
### 5.2 Channel
Channel是Go语言中用于通信的机制，Channel是Go语言中的内置类型。Channel可以传递任意类型的数据，Channel可以进行读写，也可以进行广播，也可以进行关闭。Channel可以进行阻塞，也可以进行非阻塞。

Channel主要用于协程之间的通信，比如主线程创建了10个协cheng线程，每个协cheng线程都向主线程发送消息，那么主线程就可以使用Channel进行通信，接收到消息后，就可以处理。

1. 无缓存Channel声明与使用

```go
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
```
程序解释: 创建无缓存channel，并使用go routine发送数据，然后使用for range循环接收数据。

2. 带缓存Chanel的声明与使用

无消费者
```go
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
```
生产消费者模式：
```go
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
```
### 5.3 select
select语句用于处理多个channel，select会一直等待，直到有channel有数据可读，然后执行对应的case。

无超时的select:
```go
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
```
带超时的select:
```go
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
```
select会阻塞，直到有case可以执行，如果超时，则执行case后面的代码。

循环监听多个channel

```go
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
```

- 带有退出信号的监听

```go
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
```
- 关闭channel后的读取
```go
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
```
- context与channel判断timeout
```go
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
```
-context与channel判断cancel
```go
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
```