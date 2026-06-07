package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rpcURL := "https://mainnet.infura.io/v3/9edc71dd25e5412b8f973b8651981df3"
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Println("连接节点失败:", err)
		return
	}
	txHash := common.HexToHash("0x3b626a01fe7e631ac42cd4c749196e5dfc40726c92c546810083ec708893d896")
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		fmt.Println("查询交易失败:", err)
		return
	}
	// To() / Value() / Gas() / GasPrice() / Data() / Nonce()
	fmt.Println("=== Transaction ===")
	fmt.Printf("Tx Hash       : %s\n", tx.Hash().Hex())
	fmt.Printf("Tx Pending    : %t\n", isPending) // 判断交易是否还在等待打包
	fmt.Printf("Tx To         : %v\n", tx.To())
	fmt.Printf("Tx Value      : %v\n", tx.Value())
	fmt.Printf("Tx Gas        : %v\n", tx.Gas())
	fmt.Printf("Tx GasPrice   : %v\n", tx.GasPrice())
	fmt.Printf("Tx Data       : %v\n", tx.Data())
	fmt.Printf("Tx Nonce      : %v\n", tx.Nonce()) // 账户已发出的交易数
	fmt.Println("==========================")

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Receipt ===")
	fmt.Printf("Tx Hash      : %s\n", receipt.TxHash.Hex())
	fmt.Printf("Block Hash   : %s\n", receipt.BlockHash.Hex())
	fmt.Printf("Block Number : %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Status       : %d\n", receipt.Status)
	fmt.Printf("Logs         : %v\n", receipt.Logs)
	fmt.Println("==========================")
}
