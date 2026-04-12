package main

import (
	"errors"
	"fmt"
)

/*
* @Description:
完整的计算器实现
包含加减乘除和统计功能
*/
type Recorder struct {
	ID          uint32
	Description string
}
type History struct {
	recorder []Recorder
}

func (h *History) addOperation(operation string) {
	_len := len(h.recorder)
	current_id := uint32(_len) + 1
	current_recorder := Recorder{ID: current_id, Description: operation}
	h.recorder = append(h.recorder, current_recorder)
}
func add(a, b int, history *History) int {
	history.addOperation(fmt.Sprintf("%d + %d = %d", a, b, a+b))
	return a + b
}
func subtract(a, b int, history *History) int {
	history.addOperation(fmt.Sprintf("%d - %d = %d", a, b, a-b))
	return a - b
}
func multiply(a, b int, history *History) int {
	history.addOperation(fmt.Sprintf("%d * %d = %d", a, b, a*b))
	return a * b
}
func divide(a, b int, history *History) (int, error) {
	if b == 0 {
		history.addOperation(fmt.Sprintf("%d / %d = err", a, b))

		return 0, errors.New("除数不能为0")
	}
	history.addOperation(fmt.Sprintf("%d / %d = %d", a, b, a/b))

	return a / b, nil
}
func menu() {
	fmt.Println("1. 加")
	fmt.Println("2. 减")
	fmt.Println("3. 乘")
	fmt.Println("4. 除")
	fmt.Println("5. 统计")
	fmt.Println("6. 退出")
}
func main() {
	menu()
	history := History{recorder: []Recorder{}}
	for {
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			var a, b int
			fmt.Print("请输入两个整数: ")
			fmt.Scanln(&a, &b)
			result := add(a, b, &history)
			fmt.Printf("结果: %d\n", result)
		case 2:
			var a, b int
			fmt.Print("请输入两个整数: ")
			fmt.Scanln(&a, &b)
			result := subtract(a, b, &history)
			fmt.Printf("结果: %d\n", result)
		case 3:
			var a, b int
			fmt.Print("请输入两个整数: ")
			fmt.Scanln(&a, &b)
			result := multiply(a, b, &history)
			fmt.Printf("结果: %d\n", result)
		case 4:
			var a, b int
			fmt.Print("请输入两个整数: ")
			fmt.Scanln(&a, &b)
			result, err := divide(a, b, &history)
			if err != nil {
				fmt.Println("错误:", err)
			} else {
				fmt.Printf("结果: %d\n", result)
			}
		case 5:
			fmt.Println("操作记录:")
			for _, record := range history.recorder {
				fmt.Printf("%d: %s\n", record.ID, record.Description)
			}
		case 6:
			fmt.Println("退出程序")
			return
		default:
			fmt.Println("无效的选择，请重新输入")
			menu()
		}
	}
}
