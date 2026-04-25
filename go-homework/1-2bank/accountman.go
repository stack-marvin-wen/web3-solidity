package main

import (
	"fmt"
	"time"
)

func printAccountInfo(account Account) {
	fmt.Printf("账户ID: %d\n", account.Id)
	fmt.Printf("账户号码: %s\n", account.AccountNumber)
	fmt.Printf("账户持有者: %s\n", account.Owner)
	fmt.Printf("账户余额: %.2f\n", account.Balance)
}

func accountHeader() {
	fmt.Printf("%-20s %-20s %-20s\n", "账户号码", "账户持有者", "账户余额")
}
func accountInfoLine(account Account) {
	fmt.Printf("%-20s %-20s %-20.2f\n", account.AccountNumber, account.Owner, account.Balance)
}
func historyHeader() {
	fmt.Printf("%-5s %-20s %-10s %-15s %-20s\n", "ID", "账户号码", "金额", "交易类型", "时间戳")
}
func historyInfoLine(history TranHistory) {
	fmt.Printf("%-5d %-20s %-10.2f %-15s %-20s\n",
		history.Id, history.AccountId, history.Amount, history.Type, history.TimeStamp)
}
func listAllAccount(bk *Bank) {
	accountHeader()
	for _, account := range bk.Accounts {
		accountInfoLine(account)
	}

}
func listAllHistory(bk *Bank) {
	historyHeader()
	for _, history := range bk.TranHistories {
		historyInfoLine(history)
	}

}

func createAccount(bk *Bank) {
	fmt.Println("请输入账户持有者：")
	var owner string
	fmt.Scanln(&owner)
	newAccount := Account{
		Id:            bk.nextId,
		AccountNumber: fmt.Sprintf("%v", bk.nextId),
		Owner:         owner,
		Balance:       0.0,
	}
	bk.mu.Lock()
	bk.nextId++
	bk.mu.Unlock()

	bk.AccountList.mu.Lock()
	bk.AccountList.Accounts = append(bk.AccountList.Accounts, newAccount)
	bk.AccountList.mu.Unlock()
	fmt.Printf("账户创建成功！账户号码：%s\n", newAccount.AccountNumber)
	printAccountInfo(newAccount)
}
func deleteAccount(bk *Bank) {
	listAllAccount(bk)
	fmt.Println("请输入账户号码：")
	var accountNumber string
	fmt.Scanln(&accountNumber)
	bk.AccountList.mu.Lock()
	defer bk.AccountList.mu.Unlock()
	for i, account := range bk.AccountList.Accounts {
		if account.AccountNumber == accountNumber {
			fmt.Printf("账户信息：\n")
			accountHeader()
			accountInfoLine(account)
			bk.AccountList.Accounts = append(bk.AccountList.Accounts[:i], bk.AccountList.Accounts[i+1:]...)
			fmt.Println("\n账户销户成功！")
			return
		}
	}
	fmt.Println("未找到账户信息！")
}
func queryAccountBalance(bk *Bank) {
	listAllAccount(bk)
	fmt.Println("请输入账户号码：")
	var accountNumber string
	fmt.Scanln(&accountNumber)
	bk.AccountList.mu.Lock()
	defer bk.AccountList.mu.Unlock()
	for _, account := range bk.AccountList.Accounts {
		if account.AccountNumber == accountNumber {
			fmt.Printf("账户余额：%.2f\n", account.Balance)
			return
		}
	}
	fmt.Println("未找到账户信息！")
}

func storeBalance(bk *Bank, accountNumber string) {
	var amount float64
	fmt.Println("请输入存款金额：")
	fmt.Scanln(&amount)
	if amount <= 0 {
		fmt.Println("存款金额必须大于0！")
		return
	}
	bk.AccountList.mu.Lock()
	defer bk.AccountList.mu.Unlock()
	for i, account := range bk.AccountList.Accounts {
		if account.AccountNumber == accountNumber {
			bk.AccountList.Accounts[i].Balance += amount
			fmt.Printf("存款成功！当前余额：%.2f\n", bk.AccountList.Accounts[i].Balance)
			history := TranHistory{
				Id:          bk.nextHisId,
				AccountId:   account.AccountNumber,
				Amount:      amount,
				Type:        "deposit",
				SentAccount: "",
				TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
			}
			bk.TranHistories = append(bk.TranHistories, history)
			return
		}
	}
	fmt.Println("未找到账户信息！")
}
func withdraw(bk *Bank, accountNumber string) {
	var amount float64
	fmt.Println("请输入取款金额：")
	fmt.Scanln(&amount)
	if amount <= 0 {
		fmt.Println("取款金额必须大于0！")
		return
	}
	bk.AccountList.mu.Lock()
	defer bk.AccountList.mu.Unlock()
	for i, account := range bk.AccountList.Accounts {
		if account.AccountNumber == accountNumber {
			if bk.AccountList.Accounts[i].Balance < amount {
				fmt.Println("余额不足，无法取款！")
				return
			}
			bk.AccountList.Accounts[i].Balance -= amount
			fmt.Printf("取款成功！当前余额：%.2f\n", bk.AccountList.Accounts[i].Balance)
			history := TranHistory{
				Id:          bk.nextHisId,
				AccountId:   account.AccountNumber,
				Amount:      amount,
				Type:        "withdrawal",
				SentAccount: "",
				TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
			}
			bk.TranHistories = append(bk.TranHistories, history)
			return
		}
	}
	fmt.Println("未找到账户信息！")
}

func transferAcc(bk *Bank, accountNumber string) {
	var sentAccount string
	var amount float64
	fmt.Println("请输入转账目标账户号码：")
	fmt.Scanln(&sentAccount)
	fmt.Println("请输入转账金额：")
	fmt.Scanln(&amount)
	if amount <= 0 {
		fmt.Println("转账金额必须大于0！")
		return
	}
	bk.AccountList.mu.Lock()
	defer bk.AccountList.mu.Unlock()
	var senderIndex, receiverIndex int = -1, -1
	for i, account := range bk.AccountList.Accounts {
		if account.AccountNumber == accountNumber {
			senderIndex = i
		}
		if account.AccountNumber == sentAccount {
			receiverIndex = i
		}
	}
	if senderIndex == -1 || receiverIndex == -1 {
		fmt.Println("未找到账户信息！")
		return
	}
	if bk.AccountList.Accounts[senderIndex].Balance < amount {
		fmt.Println("余额不足，无法转账！")
		return
	}
	bk.AccountList.Accounts[senderIndex].Balance -= amount
	bk.AccountList.Accounts[receiverIndex].Balance += amount
	fmt.Printf("转账成功！当前余额：%.2f\n", bk.AccountList.Accounts[senderIndex].Balance)
	history := TranHistory{
		Id:          bk.nextHisId,
		AccountId:   bk.AccountList.Accounts[senderIndex].AccountNumber,
		Amount:      amount,
		Type:        "transfer",
		SentAccount: bk.AccountList.Accounts[receiverIndex].AccountNumber,
		TimeStamp:   time.Now().Format("2006-01-02 15:04:05"),
	}
	bk.TranHistories = append(bk.TranHistories, history)
}
