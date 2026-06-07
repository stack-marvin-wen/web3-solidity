package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, cannel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cannel()

	rpcURL := "https://mainnet.infura.io/v3/9edc71dd25e5412b8f973b8651981df3"

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Println("连接节点失败:", err)
		return
	}
	// 区块操作，查询指定区块
	fmt.Println("1. 区块操作，查询指定区块")
	blockNum := big.NewInt(12345678)
	block, err := client.BlockByNumber(ctx, blockNum)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Block ===")
	fmt.Printf("Block Number  : %d\n", block.Number().Uint64())
	fmt.Printf("Block Hash    : %s\n", block.Hash().Hex())
	fmt.Printf("Block Time    : %s\n", time.Unix(int64(block.Time()), 0).Format(time.RFC3339))
	fmt.Println("==========================")
	fmt.Println("2. 控制请求频率")
	rateLimit := 1000 * time.Millisecond // 每两百毫秒请求一次
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()
	start, end := 170000, 180000
	// for i := start; i < end; i++ {
	// 	<-ticker.C
	// 	block, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Block Number  : %d\n", block.Number().Uint64())
	// 	fmt.Printf("Block Hash    : %s\n", block.Hash().Hex())
	// }
	fmt.Println("3. 指数退避算法请求")
	for i := start; i < end; i++ {
		block, err := batchBlockReqRetry(ctx, client, big.NewInt(int64(i)), 10)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Block Number  : %d\n", block.Number().Uint64())
		fmt.Printf("Block Hash    : %s\n", block.Hash().Hex())
	}

}
func batchBlockReqRetry(ctx context.Context, client *ethclient.Client, blockNum *big.Int, retrynum int) (*types.Block, error) {
	for i := 0; i < retrynum; i++ {
		ctx, cannel := context.WithTimeout(ctx, 5*time.Second)
		defer cannel()
		block, err := client.BlockByNumber(ctx, blockNum)
		if err == nil {
			return block, err
		}
		if i < retrynum-1 {
			backoff := time.Duration(1<<i) * 500 * time.Millisecond
			time.Sleep(backoff)
		}
	}
	return nil, fmt.Errorf("batchBlockReqRetry failed")
}
