package main

import "fmt"

func accountManagement(bk *Bank) {
	m := Menu{}
	for {
		m.accountManagementMenu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			createAccount(bk)
		case 2:
			deleteAccount(bk)
		case 3:
			queryAccountBalance(bk)
		case 4:
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}

func tranManagement(bk *Bank) {
	m := Menu{}
	var accountNumber string
	fmt.Println("请输入操作账户号码：")
	fmt.Scanln(&accountNumber)
	for {
		m.transactionManagementMenu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			storeBalance(bk, accountNumber)
			bk.nextHisId++
		case 2:
			withdraw(bk, accountNumber)
			bk.nextHisId++
		case 3:
			transferAcc(bk, accountNumber)
			bk.nextHisId++
		case 4:
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}
func queryService(bk *Bank) {
	m := Menu{}
	for {
		m.queryServiceMenu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			listAllAccount(bk)
		case 2:
			listAllHistory(bk)
		case 3:
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}
func main() {
	account := []Account{
		{Id: 1, AccountNumber: "1234567890", Owner: "Alice", Balance: 1000.0},
		{Id: 2, AccountNumber: "0987654321", Owner: "Bob", Balance: 2000.0},
	}
	tranHistory := []TranHistory{
		{Id: 1, AccountId: "1", Amount: 500.0, SentAccount: "1234567890", Type: "deposit"},
		{Id: 2, AccountId: "2", Amount: 300.0, SentAccount: "0987654321", Type: "withdrawal"},
	}
	bk := Bank{
		AccountList: AccountList{
			Accounts: account,
		},
		TranHistoriesList: TranHistoriesList{
			TranHistories: tranHistory,
		},
		nextId: 647291520420003123,
	}
	m := Menu{}
	for {
		m.menu()
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			accountManagement(&bk)
		case 2:
			tranManagement(&bk)
		case 3:
			queryService(&bk)
		case 4:
			fmt.Println("退出系统")
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}
