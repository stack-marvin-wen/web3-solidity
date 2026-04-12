package main

import "fmt"

// 基本函数
func add(a int, b int) int {
	fmt.Println("基本函数:完成a+b")
	return a + b
}

// 类型简写
func multiply(a, b int) int {
	fmt.Println("类型简写:完成a*b")
	return a * b
}

// 多返回值
func divide(a, b int) (int, int) {
	defer func() (int, int) {
		if r := recover(); r != nil {
			fmt.Println("panic被恢复:", r)
			return 0, 0
		}
		return a / b, a % b
	}()

	if b == 0 {
		panic("除数不能为0")
	}
	fmt.Println("多返回值:完成a/b")
	return a / b, a % b
}

// 命名返回值
func calculate(a, b int) (sum int, product int) {
	fmt.Println("命名返回值:开始计算")
	sum = a + b
	product = a * b
	fmt.Println("命名返回值:完成计算")
	return // 直接return会返回sum和product的值
}
func sum(nums ...int) int {
	fmt.Println("可变参数:计算所有参数的和")
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}
func main() {
	fmt.Println(add(1, 2))
	fmt.Println(multiply(3, 2))
	fmt.Println(divide(10, 0))
	fmt.Println(calculate(3, 2))
	fmt.Println(sum(1, 2, 3, 4, 5))
	fmt.Println(sum(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
}
