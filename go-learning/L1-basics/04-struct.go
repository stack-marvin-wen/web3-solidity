package main

import "fmt"

// 定义结构体
type Person struct {
	Name string
	Age  int
}

// 结构体方法
func (p Person) GetInfo() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

// 值接收者方法（不能修改结构体）
func (p Person) IncrementAgeWrong() {
	p.Age++ // 这不会修改原始结构体
}

// 指针接收者方法（可以修改结构体）
func (p *Person) IncrementAge() {
	p.Age++ // 这会修改原始结构体
}

// 结构体嵌套
type Employee struct {
	Person // 嵌套Person结构体
	ID     string
}

func (e Employee) GetEmployeeInfo() string {
	return fmt.Sprintf("ID: %s, %s", e.ID, e.GetInfo())
}

// 接口
type Speaker interface {
	Speak() string
}

func (p Person) Speak() string {
	return "I am a speaker"
}
func main() {
	fmt.Println("=== 结构体示例 ===")
	// 创建结构体实例
	p1 := Person{Name: "Alice", Age: 30}
	fmt.Println(p1.GetInfo())

	// 尝试使用值接收者方法修改年龄
	p1.IncrementAgeWrong()
	fmt.Println("After IncrementAgeWrong:", p1.GetInfo()) // 年龄不会改变
	// 简短初始化
	p2 := Person{"Bob", 25}
	fmt.Println(p2.GetInfo())

	// 使用指针接收者方法修改年龄
	p2.IncrementAge()
	fmt.Println("After IncrementAge:", p2.GetInfo()) // 年龄会改变
	// 部分初始化
	p3 := Person{Age: 35}
	fmt.Println(p3.GetInfo())
	// 嵌套结构体
	fmt.Println("=== 结构体嵌套示例 ===")
	emp := Employee{
		Person: Person{Name: "Charlie", Age: 28},
		ID:     "E123",
	}
	fmt.Println(emp.GetEmployeeInfo())
	fmt.Println(emp.Person.GetInfo())
	fmt.Println(emp.GetInfo())
	fmt.Println(emp.Name)
}
