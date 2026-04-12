package main

import (
	"fmt"
)

/**
 * 02-control-flow.go
 * if 语句
 * switch
 * for 循环
 * defer 语句
 * panic 和 recover
 */
func if_statement() {
	fmt.Println("=== if/else示例 ===")
	age := 20
	if age >= 18 {
		fmt.Println("你已成年")
	} else {
		fmt.Println("你还未成年")
	}
	// if中支持初始化变量
	if score := 85; score >= 90 {
		fmt.Println("成绩优秀")
	}
	if score := 85; score >= 90 {
		fmt.Println("成绩优秀")
	} else {
		fmt.Println("成绩一般")
	}
}
func switch_statement() {
	fmt.Println("\n=== switch示例 ===")
	day := 3
	// switch中默认是不穿透的，如果需要穿透需要使用fallthrough关键字
	switch day {
	case 1:
		fmt.Println("今天是周一")
	case 2:
		fmt.Println("今天是周二")
	case 3:
		fmt.Println("今天是周三")
	case 4:
		fmt.Println("今天是周四")
	case 5:
		fmt.Println("今天是周五")
	case 6:
		fmt.Println("今天是周六")
	case 7:
		fmt.Println("今天是周日")
	default:
		fmt.Println("未知的日期")
	}
	// switch中多条件判断
	switch day {
	case 1, 2, 3, 4, 5:
		fmt.Println("今天是工作日")
	case 6, 7:
		fmt.Println("今天是周末")
	default:
		fmt.Println("未知的日期")
	}
	// switch中条件表达式
	score := 85
	switch {
	case score >= 90:
		fmt.Println("成绩优秀")
	case score >= 60:
		fmt.Println("成绩及格")
	default:
		fmt.Println("成绩不及格")
	}

}
func for_loop() {
	// 基本结构
	fmt.Println("\n=== for循环示例 ===")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
	// 类似while循环
	j := 0
	for j < 5 {
		fmt.Println(j)
		j++
	}
	// 无限循环
	// for {
	// 	fmt.Println("无限循环")
	// }

	// map遍历
	fmt.Println("\n=== map遍历 ===")
	m := map[string]int{
		"apple":  12,
		"orange": 15,
	}
	for key, value := range m {
		fmt.Printf("%s: %d\n", key, value)
	}
	// slice遍历
	fmt.Println("\n=== slice遍历 ===")
	s := []int{1, 2, 3, 4, 5}
	for index, value := range s {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}
	// 数组遍历
	fmt.Println("\n=== 数组遍历 ===")
	arr := [5]int{1, 2, 3, 4, 5}
	for index, value := range arr {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}
}
func defer_statement() {
	fmt.Println("\n=== defer示例 ===")
	// 基本defer
	fmt.Println("1. 基本defer,延迟执行到函数返回前:")
	defer fmt.Println("defer语句被执行了")
	fmt.Println("函数正在执行中...")
	// 多个defer执行顺序(LIFO)
	defer fmt.Println("defer语句2")
	defer fmt.Println("defer语句3")
	fmt.Println("函数执行结束")
	// defer在return之后执行
	fmt.Println("2. defer在return之后执行:")
	fmt.Println("函数正在执行中...")
	fmt.Println(returnWithDefer())
	// defer在捕获变量的时机
	deferValueCapture()
	// defer在闭包捕获最终值
	deferClosureDemo()
	// defer在panic之后也会执行
	deferPanicDemo()
	// defer用于资源清理
	deferResourceCleanup()
}
func deferResourceCleanup() {
	fmt.Println("\n=== defer用于资源清理示例 ===")
	// 模拟打开一个资源
	if err := mockReadFile("02-control-flow.go"); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}
	if err := mockReadFile("nonexistent.txt"); err != nil {
		fmt.Printf("  错误: %v\n", err)
	}
}

// mockReadFile 模拟文件读取，演示defer资源清理
func mockReadFile(filename string) error {
	fmt.Printf("  打开文件: %s\n", filename)
	defer fmt.Printf("  关闭文件: %s\n", filename)
	if filename == "nonexistent.txt" {
		return fmt.Errorf("文件不存在")
	}
	fmt.Printf("  读取文件内容: %s\n", filename)
	return nil
}
func deferPanicDemo() {
	fmt.Println("\n=== defer在panic之后执行示例 ===")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("defer捕获到panic: %v\n", r)
		}
	}()
	panic("发生了一个错误")
}
func deferClosureDemo() {
	fmt.Println("\n=== defer闭包捕获变量示例 ===")
	x := 10
	defer func() {
		fmt.Printf("defer捕获的x值: %d\n", x)
	}()
	x++
	fmt.Printf("修改后的x值: %d\n", x)
}
func deferValueCapture() {
	fmt.Println("\n=== defer捕获变量示例 ===")
	x := 10
	x++
	defer fmt.Printf("defer捕获的x值: %d\n", x)
	x++
	fmt.Printf("修改后的x值: %d\n", x)
}
func returnWithDefer() int {
	defer fmt.Println("  defer: 在return之后执行")
	fmt.Println("  return: 先准备返回值")
	return 42
}

func panic_statement() {
	fmt.Println("\n=== panic示例 ===")
	panic_recover()
	panic_norecover()
}
func panic_norecover() {
	fmt.Println("\n=== panic不恢复示例 ===")
	panic("发生一个错误")
	fmt.Println("这行代码不执行")
}
func panic_recover() {
	fmt.Println("\n=== panic恢复示例 ===")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic被恢复:", r)
		}
	}()
	panic("发生一个错误")
}
func main() {
	if_statement()
	switch_statement()
	for_loop()
	defer_statement()
	panic_statement()
}
