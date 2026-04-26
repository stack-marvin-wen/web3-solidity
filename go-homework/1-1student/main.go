package main

import (
	"fmt"
	"os"
)

func menu() {
	/*
		增加学生
		删除学生
		修改学生信息
		查询学生信息
		列出所有学生信息
		退出系统
	*/
	fmt.Println("欢迎来到学生管理系统")
	fmt.Println("1. 增加学生")
	fmt.Println("2. 删除学生")
	fmt.Println("3. 修改学生信息")
	fmt.Println("4. 查询学生信息")
	fmt.Println("5. 列出所有学生信息")
	fmt.Println("6. 退出系统")
	fmt.Println("请输入您的选择：")
}

/*
学生id
学生姓名
学生学号
学生年龄
学生性别
学生班级
学生电话
学生邮箱
学生描述
学生地址
*/
func printHeader() {
	header := []string{
		"ID",
		"name",
		"age",
		"sex",
		"class",
		"tel",
		"email",
		"desc",
		"add",
	}
	fmt.Printf("%5s %12s %5s %12s %24s %15s\n",
		header[0], header[1], header[2], header[3], header[4],
		header[5])
}
func printStudent(student Student) {
	fmt.Printf("%5d %12s %5d %12s %24s %15s\n",
		student.Id, student.Name, student.Age, student.Grander, student.Class,
		student.Phone)
}
func updateStudentInfo(student *Student) {
	fmt.Println("请输入修改信息：")
	fmt.Print("姓名：")
	fmt.Scanln(&student.Name)
	fmt.Print("年龄：")
	fmt.Scanln(&student.Age)
	fmt.Print("性别：")
	fmt.Scanln(&student.Grander)
	fmt.Print("班级：")
	fmt.Scanln(&student.Class)
	fmt.Print("电话：")
	fmt.Scanln(&student.Phone)
	fmt.Print("邮箱：")
	fmt.Scanln(&student.Email)
	fmt.Print("描述：")
	fmt.Scanln(&student.Description)
	fmt.Print("地址：")
	fmt.Scanln(&student.Address)
	fmt.Println("学生信息已更新！")

}
func addStudent(students []Student, id *int) []Student {
	var student Student
	fmt.Println("请输入学生信息：")
	student.Id = *id
	fmt.Print("姓名：")
	fmt.Scanln(&student.Name)
	fmt.Print("年龄：")
	fmt.Scanln(&student.Age)
	fmt.Print("性别：")
	fmt.Scanln(&student.Grander)
	fmt.Print("班级：")
	fmt.Scanln(&student.Class)
	fmt.Print("电话：")
	fmt.Scanln(&student.Phone)
	fmt.Print("邮箱：")
	fmt.Scanln(&student.Email)
	fmt.Print("描述：")
	fmt.Scanln(&student.Description)
	fmt.Print("地址：")
	fmt.Scanln(&student.Address)
	fmt.Println("学生信息已保存！")
	students = append(students, student)
	*id++
	return students
}
func deleteStudent(students []Student) {
	var id int
	fmt.Println("请输入学生ID：")
	fmt.Scanln(&id)
	for i, student := range students {
		if student.Id == id {
			fmt.Println("删除学生信息：")
			printHeader()
			printStudent(student)
			students = append(students[:i], students[i+1:]...)
			return
		}
	}
	fmt.Println("未找到学生信息！")
}
func listStudents(students []Student) {
	fmt.Println("所有学生信息：")
	printHeader()
	for _, student := range students {
		printStudent(student)
	}
}
func updateStudent(students []Student) {
	var id int
	fmt.Println("请输入学生ID：")
	fmt.Scanln(&id)
	for _, student := range students {
		if student.Id == id {
			fmt.Println("=======当前学生信息:=======")
			printHeader()
			printStudent(student)
			updateStudentInfo(&student)
			return
		}
	}
}
func queryStudent(students []Student) {
	var id int
	fmt.Println("请输入学生ID：")
	fmt.Scanln(&id)
	for _, student := range students {
		if student.Id == id {
			fmt.Println("查询学生信息：")
			printHeader()
			printStudent(student)
			return
		}
	}
	fmt.Println("未找到学生信息！")
}
func main() {
	students := []Student{
		{Id: 0, Name: "张三", Age: 20, Grander: "男", Class: "计算机科学与技术", Phone: "1234567890", Email: "zhangsan@example.com", Address: "中国", Description: "无"},
		{Id: 1, Name: "李四", Age: 18, Grander: "女", Class: "计算机科学与技术", Phone: "1234567890", Email: "lisi@example.com", Address: "中国", Description: "无"},
		{Id: 2, Name: "王五", Age: 19, Grander: "男", Class: "计算机科学与技术", Phone: "1234567890", Email: "wangwu@example.com", Address: "中国", Description: "无"},
		{Id: 3, Name: "赵六", Age: 20, Grander: "女", Class: "计算机科学与技术", Phone: "1234567890", Email: "zhaoliu@example.com", Address: "中国", Description: "无"},
		{Id: 4, Name: "孙七", Age: 21, Grander: "男", Class: "计算机科学与技术", Phone: "1234567890", Email: "sunqi@example.com", Address: "中国", Description: "无"},
		{Id: 5, Name: "周八", Age: 22, Grander: "女", Class: "计算机科学与技术", Phone: "1234567890", Email: "zhouba@example.com", Address: "中国", Description: "无"},
		{Id: 6, Name: "吴九", Age: 23, Grander: "男", Class: "计算机科学与技术", Phone: "1234567890", Email: "wujing@example.com", Address: "中国", Description: "无"},
		{Id: 7, Name: "郑十", Age: 24, Grander: "女", Class: "计算机科学与技术", Phone: "1234567890", Email: "zhengshi@example.com", Address: "中国", Description: "无"},
	}
	nextId := 8
	for {
		menu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			students = addStudent(students, &nextId)
		case 2:
			deleteStudent(students)
		case 3:
			updateStudent(students)
		case 4:
			queryStudent(students)
		case 5:
			listStudents(students)
		case 6:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("无效的选择，请重新输入")
			continue
		}
	}
}
