package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcIRL := "wss://mainnet.infura.io/v3/9edc71dd25e5412b8f973b8651981df3"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 连接以太坊节点
	client, err := ethclient.DialContext(ctx, rpcIRL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	headerCh := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) // 监听中断信号
	for {
		select {
		case h := <-headerCh:
			fmt.Printf("[%s] New Block - Number: %d, Hash: %s\n",
				time.Now().Format(time.RFC3339),
				h.Number.Uint64(),
				h.Hash().Hex(),
			)
		case err := <-sub.Err():
			panic(err)
		case <-ctx.Done():
			fmt.Println("Context timeout, exiting...")
		case <-sigCh:
			fmt.Println("Interrupt signal, exiting...")
			return
		}
	}
}
