package main

import (
	"sync"
)

type Account struct {
	Id            int64   `json:"id"`
	AccountNumber string  `json:"account_number"`
	Owner         string  `json:"owner"`
	Balance       float64 `json:"balance"`
}
type TranHistory struct {
	Id          int64   `json:"id"`
	AccountId   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	SentAccount string  `json:"sent_account"`
	Type        string  `json:"type"` // "deposit", "withdrawal", "transfer"
	TimeStamp   string  `json:"timestamp"`
}

type AccountList struct {
	Accounts []Account `json:"accounts"`
	mu       sync.Mutex
}
type TranHistoriesList struct {
	TranHistories []TranHistory `json:"tran_histories"`
	mu            sync.Mutex
}
type Bank struct {
	AccountList
	TranHistoriesList
	nextId    int64
	nextHisId int64
	mu        sync.Mutex
}
