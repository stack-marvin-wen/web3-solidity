package block

import (
	"NFTAuctionServer/utils"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type ContractAuction struct {
	Seller        common.Address
	NFTContract   common.Address
	TokenId       *big.Int
	PaymentToken  common.Address
	StartPriceUsd *big.Int
	HighestBidAmt *big.Int
	HighestBidUsd *big.Int
	HighestBidder common.Address
	EndTime       *big.Int
	Active        bool
}

type ContractAuctionCreate struct {
	NFTContract   common.Address
	TokenId       *big.Int
	PaymentToken  common.Address
	StartPriceUsd *big.Int
	DurationHours *big.Int
}

type NFTMintRequest struct {
	URI string
}

type NFTApproveRequest struct {
	TokenID *big.Int
	Spender common.Address
}

type BidLog struct {
	AuctionID   *big.Int
	Bidder      common.Address
	Amount      *big.Int
	AmountUsd   *big.Int
	TxHash      common.Hash
	BlockNumber uint64
	LogIndex    uint
}

type AuctionEventLog struct {
	EventType     string
	AuctionID     *big.Int
	Seller        common.Address
	NFTContract   common.Address
	TokenId       *big.Int
	PaymentToken  common.Address
	StartPriceUsd *big.Int
	EndTime       *big.Int
	Bidder        common.Address
	Winner        common.Address
	Amount        *big.Int
	AmountUsd     *big.Int
	TxHash        common.Hash
	BlockNumber   uint64
	LogIndex      uint
}

func ValidateNFTAuctionCreate(auction *ContractAuctionCreate, privateKey *ecdsa.PrivateKey) error {
	parseNFTABI, err := utils.ParseABI("abi/NFTAInstance.json")
	if err != nil {
		logger.Error("解析 NFT ABI 文件失败: ", err)
		return err
	}

	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contractAddr := auction.NFTContract
	data, err := parseNFTABI.Pack("ownerOf", auction.TokenId)
	if err != nil {
		logger.Error("打包 ownerOf 调用数据失败: ", err)
		return err
	}

	ownerOutput, err := client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: data}, nil)
	if err != nil {
		logger.Error("查询 NFT 所有权失败: ", err)
		return fmt.Errorf("NFT 不存在或 tokenId 无效: %w", err)
	}

	ownerValues, err := parseNFTABI.Unpack("ownerOf", ownerOutput)
	if err != nil || len(ownerValues) != 1 {
		logger.Error("解析 ownerOf 返回值失败: ", err)
		return fmt.Errorf("解析 NFT 所有权失败")
	}
	owner, ok := ownerValues[0].(common.Address)
	if !ok {
		return fmt.Errorf("解析 NFT 所有权失败")
	}

	caller := privateKeyToAddress(privateKey)
	if owner != caller {
		return fmt.Errorf("NFT 当前拥有者不是创建拍卖的账户")
	}
	auctionContractAddr := common.HexToAddress(utils.GetAuctionContractAddress())

	approvedData, err := parseNFTABI.Pack("getApproved", auction.TokenId)
	if err != nil {
		logger.Error("打包 getApproved 调用数据失败: ", err)
		return err
	}

	approvedOutput, err := client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: approvedData}, nil)
	if err != nil {
		logger.Error("查询 NFT 单个授权失败: ", err)
		return err
	}
	approvedValues, err := parseNFTABI.Unpack("getApproved", approvedOutput)
	if err != nil || len(approvedValues) != 1 {
		logger.Error("解析 getApproved 返回值失败: ", err)
		return fmt.Errorf("解析 NFT 授权状态失败")
	}
	approved, ok := approvedValues[0].(common.Address)
	if !ok {
		return fmt.Errorf("解析 NFT 授权状态失败")
	}
	if approved == auctionContractAddr {
		return nil
	}

	allApprovedData, err := parseNFTABI.Pack("isApprovedForAll", owner, auctionContractAddr)
	if err != nil {
		logger.Error("打包 isApprovedForAll 调用数据失败: ", err)
		return err
	}
	allApprovedOutput, err := client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: allApprovedData}, nil)
	if err != nil {
		logger.Error("查询 NFT 全局授权失败: ", err)
		return err
	}
	allApprovedValues, err := parseNFTABI.Unpack("isApprovedForAll", allApprovedOutput)
	if err != nil || len(allApprovedValues) != 1 {
		logger.Error("解析 isApprovedForAll 返回值失败: ", err)
		return fmt.Errorf("解析 NFT 授权状态失败")
	}
	allApproved, ok := allApprovedValues[0].(bool)
	if !ok {
		return fmt.Errorf("解析 NFT 授权状态失败")
	}
	if !allApproved {
		return fmt.Errorf("NFT 还没有授权给拍卖合约")
	}

	return nil
}

func MintNFT(uri string) (string, error) {
	parseNFTABI, err := utils.ParseABI("abi/NFTAInstance.json")
	if err != nil {
		logger.Error("解析 NFT ABI 文件失败: ", err)
		return "", err
	}

	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		logger.Error("获取链 ID 失败: ", err)
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(utils.GetPrivateKey(), "0x"))
	if err != nil {
		logger.Error("解析私钥失败: ", err)
		return "", err
	}

	contractAddr := common.HexToAddress(utils.GetNFTContractAddress())
	data, err := parseNFTABI.Pack("mint", uri)
	if err != nil {
		logger.Error("打包 mint 调用数据失败: ", err)
		return "", err
	}

	nonce, err := client.PendingNonceAt(ctx, privateKeyToAddress(privateKey))
	if err != nil {
		logger.Error("获取 nonce 失败: ", err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("获取 gas price 失败: ", err)
		return "", err
	}

	value := big.NewInt(0)
	value.SetString("50000000000000000", 10)
	tx := types.NewTransaction(nonce, contractAddr, value, 5_000_000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		logger.Error("签名 mint 交易失败: ", err)
		return "", err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		logger.Error("发送 mint 交易失败: ", err)
		return "", err
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		logger.Error("等待 mint 交易回执失败: ", err)
		return "", err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.Error("mint 交易执行失败, receipt status: ", receipt.Status)
		return "", fmt.Errorf("mint 交易执行失败")
	}

	return signedTx.Hash().Hex(), nil
}

func ApproveNFT(tokenID *big.Int) (string, error) {
	parseNFTABI, err := utils.ParseABI("abi/NFTAInstance.json")
	if err != nil {
		logger.Error("解析 NFT ABI 文件失败: ", err)
		return "", err
	}

	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		logger.Error("获取链 ID 失败: ", err)
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(utils.GetPrivateKey(), "0x"))
	if err != nil {
		logger.Error("解析私钥失败: ", err)
		return "", err
	}

	contractAddr := common.HexToAddress(utils.GetNFTContractAddress())
	spender := common.HexToAddress(utils.GetAuctionContractAddress())
	data, err := parseNFTABI.Pack("approve", spender, tokenID)
	if err != nil {
		logger.Error("打包 approve 调用数据失败: ", err)
		return "", err
	}

	nonce, err := client.PendingNonceAt(ctx, privateKeyToAddress(privateKey))
	if err != nil {
		logger.Error("获取 nonce 失败: ", err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("获取 gas price 失败: ", err)
		return "", err
	}

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), 300_000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		logger.Error("签名 approve 交易失败: ", err)
		return "", err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		logger.Error("发送 approve 交易失败: ", err)
		return "", err
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		logger.Error("等待 approve 交易回执失败: ", err)
		return "", err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.Error("approve 交易执行失败, receipt status: ", receipt.Status)
		return "", fmt.Errorf("approve 交易执行失败")
	}

	return signedTx.Hash().Hex(), nil
}

/**
 * 获取 NFT 拍卖列表
 */
func GetNFTAuctionList() ([]ContractAuction, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		logger.Error("解析 ABI 文件失败: ", err)
		return nil, err
	}
	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return nil, err
	}
	defer client.Close()
	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())

	data, err := parseAuctionABI.Pack("getAuctions")
	if err != nil {
		logger.Error("打包数据失败: ", err)
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	ctx, cancel := utils.Conext()
	defer cancel()
	output, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		logger.Error("调用合约失败: ", err)
		return nil, err
	}
	var onChainAuctions []ContractAuction
	err = parseAuctionABI.UnpackIntoInterface(&onChainAuctions, "getAuctions", output)
	if err != nil {
		logger.Error("解包数据失败: ", err)
		return nil, err
	}
	return onChainAuctions, nil
}

/***
 * 获取 NFT 拍卖 ByID
 */
func GetNFTAuctionByID(id *big.Int) (ContractAuction, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		logger.Error("解析 ABI 文件失败: ", err)
		return ContractAuction{}, err
	}
	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return ContractAuction{}, err
	}
	defer client.Close()

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	data, err := parseAuctionABI.Pack("getAuction", id)
	if err != nil {
		logger.Error("打包 getAuction 调用数据失败: ", err)
		return ContractAuction{}, err
	}

	ctx, cancel := utils.Conext()
	defer cancel()
	output, err := client.CallContract(ctx, ethereum.CallMsg{To: &contractAddr, Data: data}, nil)
	if err != nil {
		logger.Error("查询单个拍卖失败: ", err)
		return ContractAuction{}, err
	}

	values, err := parseAuctionABI.Unpack("getAuction", output)
	if err != nil {
		logger.Error("解包单个拍卖失败: ", err)
		return ContractAuction{}, err
	}
	if len(values) != 10 {
		return ContractAuction{}, fmt.Errorf("解包单个拍卖失败: 返回值数量不正确")
	}

	auction := ContractAuction{}
	if seller, ok := values[0].(common.Address); ok {
		auction.Seller = seller
	}
	if nftContract, ok := values[1].(common.Address); ok {
		auction.NFTContract = nftContract
	}
	if tokenID, ok := values[2].(*big.Int); ok {
		auction.TokenId = tokenID
	}
	if paymentToken, ok := values[3].(common.Address); ok {
		auction.PaymentToken = paymentToken
	}
	if startPriceUsd, ok := values[4].(*big.Int); ok {
		auction.StartPriceUsd = startPriceUsd
	}
	if highestBidAmount, ok := values[5].(*big.Int); ok {
		auction.HighestBidAmt = highestBidAmount
	}
	if highestBidUsd, ok := values[6].(*big.Int); ok {
		auction.HighestBidUsd = highestBidUsd
	}
	if highestBidder, ok := values[7].(common.Address); ok {
		auction.HighestBidder = highestBidder
	}
	if endTime, ok := values[8].(*big.Int); ok {
		auction.EndTime = endTime
	}
	if active, ok := values[9].(bool); ok {
		auction.Active = active
	}

	return auction, nil
}

/**
 * 创建 NFT 拍卖
 * @param auction 创建参数
 */
func CreateNFTAuction(auction *ContractAuctionCreate) (string, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		logger.Error("解析 ABI 文件失败: ", err)
		return "", err
	}
	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		logger.Error("获取链 ID 失败: ", err)
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(utils.GetPrivateKey(), "0x"))
	if err != nil {
		logger.Error("解析私钥失败: ", err)
		return "", err
	}

	if err := ValidateNFTAuctionCreate(auction, privateKey); err != nil {
		logger.Error("NFT 拍卖创建前置检查失败: ", err)
		return "", err
	}

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	data, err := parseAuctionABI.Pack(
		"createAuction",
		auction.NFTContract,
		auction.TokenId,
		auction.PaymentToken,
		auction.StartPriceUsd,
		auction.DurationHours,
	)
	if err != nil {
		logger.Error("打包数据失败: ", err)
		return "", err
	}

	nonce, err := client.PendingNonceAt(ctx, privateKeyToAddress(privateKey))
	if err != nil {
		logger.Error("获取 nonce 失败: ", err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("获取 gas price 失败: ", err)
		return "", err
	}

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), 5_000_000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		logger.Error("签名交易失败: ", err)
		return "", err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		logger.Error("发送创建拍卖交易失败: ", err)
		return "", err
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		logger.Error("等待创建拍卖交易回执失败: ", err)
		return "", err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.Error("创建拍卖交易执行失败, receipt status: ", receipt.Status)
		return "", fmt.Errorf("创建拍卖交易执行失败")
	}

	return signedTx.Hash().Hex(), nil
}

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

func BidNFTAuction(auctionID *big.Int, amount *big.Int) (string, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		logger.Error("解析 ABI 文件失败: ", err)
		return "", err
	}
	client, err := utils.ConnectToBlockchain()
	if err != nil {
		logger.Error("连接区块链失败: ", err)
		return "", err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		logger.Error("获取链 ID 失败: ", err)
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(utils.GetBuyerPrivateKey(), "0x"))
	if err != nil {
		logger.Error("解析私钥失败: ", err)
		return "", err
	}

	if amount == nil || amount.Sign() <= 0 {
		return "", fmt.Errorf("bid amount 必须大于 0")
	}

	auction, err := GetNFTAuctionByID(auctionID)
	if err != nil {
		logger.Error("查询拍卖信息失败: ", err)
		return "", err
	}

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	var data []byte
	var txValue *big.Int

	if auction.PaymentToken == (common.Address{}) {
		data, err = parseAuctionABI.Pack("bidEth", auctionID)
		if err != nil {
			logger.Error("打包 bidEth 调用数据失败: ", err)
			return "", err
		}
		txValue = amount
	} else {
		data, err = parseAuctionABI.Pack("bidErc20", auctionID, amount)
		if err != nil {
			logger.Error("打包 bidErc20 调用数据失败: ", err)
			return "", err
		}
		txValue = big.NewInt(0)
	}

	nonce, err := client.PendingNonceAt(ctx, privateKeyToAddress(privateKey))
	if err != nil {
		logger.Error("获取 nonce 失败: ", err)
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("获取 gas price 失败: ", err)
		return "", err
	}

	tx := types.NewTransaction(nonce, contractAddr, txValue, 5_000_000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		logger.Error("签名出价交易失败: ", err)
		return "", err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		logger.Error("发送出价交易失败: ", err)
		return "", err
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		logger.Error("等待出价交易回执失败: ", err)
		return "", err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.Error("出价交易执行失败, receipt status: ", receipt.Status)
		return "", fmt.Errorf("出价交易执行失败")
	}

	return signedTx.Hash().Hex(), nil
}

/**
 * 结束 NFT 拍卖
 * @param auctionID 拍卖 ID
 */
func EndNFTAuction(auctionID *big.Int) error {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		return fmt.Errorf("解析 ABI 文件失败: %s ", err.Error())
	}
	client, err := utils.ConnectToBlockchain()
	if err != nil {
		return fmt.Errorf("连接区块链失败: %s ", err.Error())
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("获取链 ID 失败: %s ", err.Error())
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(utils.GetPrivateKey(), "0x"))
	if err != nil {
		return fmt.Errorf("解析私钥失败: %s ", err.Error())
	}

	auction, err := GetNFTAuctionByID(auctionID)
	if err != nil {
		return fmt.Errorf("查询拍卖信息失败: %s ", err.Error())
	}
	if !auction.Active {
		return fmt.Errorf("拍卖已结束")
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("获取当前链上时间失败: %s ", err.Error())
	}
	if auction.EndTime != nil && header.Time < auction.EndTime.Uint64() {
		return fmt.Errorf("拍卖尚未到期，请等待到期后再结束")
	}

	data, err := parseAuctionABI.Pack("endAuction", auctionID)
	if err != nil {
		return fmt.Errorf("打包 endAuction 调用数据失败: %s ", err.Error())
	}

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	nonce, err := client.PendingNonceAt(ctx, privateKeyToAddress(privateKey))
	if err != nil {
		return fmt.Errorf("获取 nonce 失败: %s ", err.Error())
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("获取 gas price 失败: %s ", err.Error())
	}

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), 5_000_000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("签名结束拍卖交易失败: %s ", err.Error())
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("发送结束拍卖交易失败: %s ", err.Error())
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return fmt.Errorf("等待结束拍卖交易回执失败 %s ", err.Error())
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("结束拍卖交易执行失败, receipt status")
	}
	return nil
}

/**
 * 出价日志
 */
func BidLogs(auctionID *big.Int) ([]BidLog, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		return nil, fmt.Errorf("解析 ABI 文件失败: %s", err.Error())
	}

	client, err := utils.ConnectToBlockchain()
	if err != nil {
		return nil, fmt.Errorf("连接区块链失败: %s", err.Error())
	}
	defer client.Close()

	event, ok := parseAuctionABI.Events["HighestBidIncreased"]
	if !ok {
		return nil, fmt.Errorf("未找到出价事件定义")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics: [][]common.Hash{
			{event.ID, common.BigToHash(auctionID)},
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询出价日志失败: %s", err.Error())
	}

	result := make([]BidLog, 0, len(logs))
	for _, item := range logs {
		values, err := parseAuctionABI.Unpack("HighestBidIncreased", item.Data)
		if err != nil {
			return nil, fmt.Errorf("解包出价日志失败: %s", err.Error())
		}
		if len(values) != 2 {
			return nil, fmt.Errorf("解包出价日志失败: 返回值数量不正确")
		}

		amount, ok := values[0].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("解包出价日志失败: amount 类型不正确")
		}
		amountUsd, ok := values[1].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("解包出价日志失败: amountUsd 类型不正确")
		}
		if len(item.Topics) < 3 {
			return nil, fmt.Errorf("解包出价日志失败: topics 数量不正确")
		}

		result = append(result, BidLog{
			AuctionID:   new(big.Int).Set(auctionID),
			Bidder:      common.BytesToAddress(item.Topics[2].Bytes()[12:]),
			Amount:      amount,
			AmountUsd:   amountUsd,
			TxHash:      item.TxHash,
			BlockNumber: item.BlockNumber,
			LogIndex:    uint(item.Index),
		})
	}

	return result, nil
}

/**
 * 获取拍卖事件日志
 */
func AuctionLogs() ([]AuctionEventLog, error) {
	parseAuctionABI, err := utils.ParseABI("abi/NFTAuction.json")
	if err != nil {
		return nil, fmt.Errorf("解析 ABI 文件失败: %s", err.Error())
	}

	client, err := utils.ConnectToBlockchain()
	if err != nil {
		return nil, fmt.Errorf("连接区块链失败: %s", err.Error())
	}
	defer client.Close()

	createdEvent, ok := parseAuctionABI.Events["AuctionCreated"]
	if !ok {
		return nil, fmt.Errorf("未找到创建拍卖事件定义")
	}
	bidEvent, ok := parseAuctionABI.Events["HighestBidIncreased"]
	if !ok {
		return nil, fmt.Errorf("未找到出价事件定义")
	}
	endedEvent, ok := parseAuctionABI.Events["AuctionEnded"]
	if !ok {
		return nil, fmt.Errorf("未找到结束拍卖事件定义")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	contractAddr := common.HexToAddress(utils.GetAuctionContractAddress())
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics:    [][]common.Hash{{createdEvent.ID, bidEvent.ID, endedEvent.ID}},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询拍卖日志失败: %s", err.Error())
	}

	result := make([]AuctionEventLog, 0, len(logs))
	for _, item := range logs {
		switch item.Topics[0] {
		case createdEvent.ID:
			values, err := parseAuctionABI.Unpack("AuctionCreated", item.Data)
			if err != nil {
				return nil, fmt.Errorf("解包创建拍卖日志失败: %s", err.Error())
			}
			if len(values) != 3 || len(item.Topics) < 4 {
				return nil, fmt.Errorf("解包创建拍卖日志失败: 返回值数量不正确")
			}
			paymentToken, ok := values[0].(common.Address)
			if !ok {
				return nil, fmt.Errorf("解包创建拍卖日志失败: paymentToken 类型不正确")
			}
			startPriceUsd, ok := values[1].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包创建拍卖日志失败: startPriceUsd 类型不正确")
			}
			endTime, ok := values[2].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包创建拍卖日志失败: endTime 类型不正确")
			}
			result = append(result, AuctionEventLog{
				EventType:     "AuctionCreated",
				Seller:        common.BytesToAddress(item.Topics[1].Bytes()[12:]),
				NFTContract:   common.BytesToAddress(item.Topics[2].Bytes()[12:]),
				TokenId:       new(big.Int).SetBytes(item.Topics[3].Bytes()),
				PaymentToken:  paymentToken,
				StartPriceUsd: startPriceUsd,
				EndTime:       endTime,
				TxHash:        item.TxHash,
				BlockNumber:   item.BlockNumber,
				LogIndex:      uint(item.Index),
			})
		case bidEvent.ID:
			values, err := parseAuctionABI.Unpack("HighestBidIncreased", item.Data)
			if err != nil {
				return nil, fmt.Errorf("解包出价日志失败: %s", err.Error())
			}
			if len(values) != 2 || len(item.Topics) < 3 {
				return nil, fmt.Errorf("解包出价日志失败: 返回值数量不正确")
			}
			amount, ok := values[0].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包出价日志失败: amount 类型不正确")
			}
			amountUsd, ok := values[1].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包出价日志失败: amountUsd 类型不正确")
			}
			result = append(result, AuctionEventLog{
				EventType:   "HighestBidIncreased",
				AuctionID:   new(big.Int).SetBytes(item.Topics[1].Bytes()),
				Bidder:      common.BytesToAddress(item.Topics[2].Bytes()[12:]),
				Amount:      amount,
				AmountUsd:   amountUsd,
				TxHash:      item.TxHash,
				BlockNumber: item.BlockNumber,
				LogIndex:    uint(item.Index),
			})
		case endedEvent.ID:
			values, err := parseAuctionABI.Unpack("AuctionEnded", item.Data)
			if err != nil {
				return nil, fmt.Errorf("解包结束拍卖日志失败: %s", err.Error())
			}
			if len(values) != 3 || len(item.Topics) < 2 {
				return nil, fmt.Errorf("解包结束拍卖日志失败: 返回值数量不正确")
			}
			winner, ok := values[0].(common.Address)
			if !ok {
				return nil, fmt.Errorf("解包结束拍卖日志失败: winner 类型不正确")
			}
			amount, ok := values[1].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包结束拍卖日志失败: amount 类型不正确")
			}
			amountUsd, ok := values[2].(*big.Int)
			if !ok {
				return nil, fmt.Errorf("解包结束拍卖日志失败: amountUsd 类型不正确")
			}
			result = append(result, AuctionEventLog{
				EventType:   "AuctionEnded",
				AuctionID:   new(big.Int).SetBytes(item.Topics[1].Bytes()),
				Winner:      winner,
				Amount:      amount,
				AmountUsd:   amountUsd,
				TxHash:      item.TxHash,
				BlockNumber: item.BlockNumber,
				LogIndex:    uint(item.Index),
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].BlockNumber == result[j].BlockNumber {
			return result[i].LogIndex < result[j].LogIndex
		}
		return result[i].BlockNumber < result[j].BlockNumber
	})

	return result, nil
}
