package main

import "fmt"

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

func main() {
	if_statement()
	switch_statement()
	for_loop()
}
