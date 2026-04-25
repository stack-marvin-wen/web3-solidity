package main

import "fmt"

type Menu struct{}

func (m *Menu) menu() {
	fmt.Println("==欢迎来到银行账户管理系统==")
	fmt.Println("1. 账户管理")
	fmt.Println("2. 交易管理")
	fmt.Println("3. 查询服务")
	fmt.Println("4. 退出")
}
func (m *Menu) accountManagementMenu() {
	fmt.Println("========账户管理==========")
	fmt.Println("1. 开户")
	fmt.Println("2. 销户")
	fmt.Println("3. 查询账户余额")
	fmt.Println("4. 退出")
}
func (m *Menu) transactionManagementMenu() {
	fmt.Println("========交易管理==========")
	fmt.Println("1. 存款")
	fmt.Println("2. 取款")
	fmt.Println("3. 转账")
	fmt.Println("4. 退出")
}
func (m *Menu) queryServiceMenu() {
	fmt.Println("========查询服务==========")
	fmt.Println("1. 查询账户信息")
	fmt.Println("2. 查询账户历史")
	fmt.Println("3. 退出")
}
