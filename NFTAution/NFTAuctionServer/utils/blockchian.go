package utils

import (
	"NFTAuctionServer/config"
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Conext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return ctx, cancel
}

/**
 * 连接区块链
 */
func ConnectToBlockchain() (*ethclient.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// 连接区块链
	rpcURL := config.Config.BlockchainConfig.RPC_URL
	// 这里可以使用 go-ethereum 的 ethclient 包来连接区块链，并使用私钥创建一个 transactor
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	return client, nil
}

/**
 * 获取私钥
 */
func GetPrivateKey() string {
	return config.Config.BlockchainConfig.Private_KEY
}

/**
 * 获取买家私钥
 */
func GetBuyerPrivateKey() string {
	return config.Config.BlockchainConfig.Buyer
}

/**
 * 获取拍卖合约地址
 */
func GetAuctionContractAddress() string {
	return config.Config.BlockchainConfig.Aution_Contract_Address
}

/**
 * 获取 NFT 合约地址
 */
func GetNFTContractAddress() string {
	return config.Config.BlockchainConfig.NFT_Contract_Address
}

/**
 * 解析 ABI 文件
 * @param path ABI 文件路径
 */
func ParseABI(path string) (abi.ABI, error) {
	abiBytes, err := os.ReadFile(path)
	if err != nil {
		return abi.ABI{}, err
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err == nil {
		return parsedABI, nil
	}

	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	if err := json.Unmarshal(abiBytes, &artifact); err != nil {
		return abi.ABI{}, err
	}
	if len(artifact.ABI) == 0 {
		return abi.ABI{}, err
	}
	return abi.JSON(strings.NewReader(string(artifact.ABI)))
}
