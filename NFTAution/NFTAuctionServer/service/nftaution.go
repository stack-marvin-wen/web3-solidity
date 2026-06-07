package service

import (
	"NFTAuctionServer/block"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type NFTAuction struct {
	ID          string `json:"id"`          // ID
	Name        string `json:"name"`        // 名称
	StartTime   string `json:"start_time"`  // 开始时间
	EndTime     string `json:"end_time"`    // 结束时间
	LowestBid   string `json:"lowest_bid"`  // 最低价
	HighestBid  string `json:"highest_bid"` // 最高价
	Description string `json:"description"` // 描述
	Status      string `json:"status"`      // 状态
	Image       string `json:"image"`       // 图片URL
	Tags        string `json:"tags"`        // 标签
}

type AuctionLog struct {
	EventType     string `json:"event_type"`
	AuctionID     string `json:"auction_id,omitempty"`
	Seller        string `json:"seller,omitempty"`
	NFTContract   string `json:"nft_contract,omitempty"`
	TokenID       string `json:"token_id,omitempty"`
	PaymentToken  string `json:"payment_token,omitempty"`
	StartPriceUsd string `json:"start_price_usd,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
	Bidder        string `json:"bidder,omitempty"`
	Winner        string `json:"winner,omitempty"`
	Amount        string `json:"amount,omitempty"`
	AmountUsd     string `json:"amount_usd,omitempty"`
	TxHash        string `json:"tx_hash"`
	BlockNumber   string `json:"block_number"`
	LogIndex      string `json:"log_index"`
}

type AuctionLogFilter struct {
	EventType   string
	AuctionID   string
	Seller      string
	NFTContract string
	TokenID     string
	Bidder      string
	Winner      string
	TxHash      string
	FromBlock   string
	ToBlock     string
}

type OverviewAution struct {
	NFTAuctionCount int `json:"nft_auction_count"` // 拍卖总数
	BidCount        int `json:"bid_count"`         // 出价总数
	BidAmount       int `json:"bid_amount"`        // 出价总金额
}

type NFTAuctionServiceInterface interface {
	GetNFTAuctionList() ([]NFTAuction, error)                          // 获取拍卖列表
	GetNFTAuctionByID(id string) (NFTAuction, error)                   // 根据ID获取拍卖详情
	CreateNFTAuction(auction *CreateNFTAuctionRequest) (string, error) // 创建拍卖
	MintNFT(uri string) (string, error)                                // 铸造NFT
	ApproveNFT(tokenID string) (string, error)                         // 授权NFT给拍卖合约
	BidNFTAuction(auctionID string, bidAmount string) (string, error)  // 拍卖出价
	GetBidLogs(id string) ([]BidLog, error)                            // 获取出价日志
	GetAuctionLogs(filter *AuctionLogFilter) ([]AuctionLog, error)     // 获取全部拍卖日志
	GetTotalNFTAuctions() (OverviewAution, error)                      // 获取拍卖总数
	EndNFTAuction(id string) error                                     // 结束拍卖
}

type CreateNFTAuctionRequest struct {
	NFTContract   string `json:"nft_contract" binding:"required"`
	TokenID       string `json:"token_id" binding:"required"`
	PaymentToken  string `json:"payment_token" binding:"required"`
	StartPriceUSD string `json:"start_price_usd" binding:"required"`
	DurationHours string `json:"duration_hours" binding:"required"`
}

type MintNFTRequest struct {
	URI string `json:"uri" binding:"required"`
}

type ApproveNFTRequest struct {
	TokenID string `json:"token_id" binding:"required"`
}

type BidNFTAuctionRequest struct {
	AuctionID string `json:"auction_id" binding:"required"`
	BidAmount string `json:"bid_amount" binding:"required"`
}

type BidLog struct {
	AuctionID   string `json:"auction_id"`
	Bidder      string `json:"bidder"`
	Amount      string `json:"amount"`
	AmountUsd   string `json:"amount_usd"`
	TxHash      string `json:"tx_hash"`
	BlockNumber string `json:"block_number"`
	LogIndex    string `json:"log_index"`
}

type EndNFTAuctionRequest struct {
	AuctionID string `json:"auction_id" binding:"required"`
}

type NFTAuctionService struct {
}

func (nft *NFTAuctionService) GetNFTAuctionList() ([]NFTAuction, error) {

	onChainAuctions, err := block.GetNFTAuctionList()
	if err != nil {
		return nil, err
	}
	auctions := make([]NFTAuction, 0, len(onChainAuctions))
	for i, auction := range onChainAuctions {
		auctions = append(auctions, NFTAuction{
			ID:         fmt.Sprintf("%d", i),
			Name:       auction.NFTContract.Hex(),
			StartTime:  "",
			EndTime:    auction.EndTime.String(),
			LowestBid:  auction.StartPriceUsd.String(),
			HighestBid: auction.HighestBidUsd.String(),
			Status:     fmt.Sprintf("%t", auction.Active),
		})
	}
	// 获取拍卖列表
	return auctions, nil
}

/**
 * 获取拍卖详情
 */
func (nft *NFTAuctionService) GetNFTAuctionByID(id string) (NFTAuction, error) {
	i_id, err := strconv.Atoi(id)
	if err != nil {
		return NFTAuction{}, err
	}
	onChainAuction, err := block.GetNFTAuctionByID(big.NewInt(int64(i_id)))
	if err != nil {
		return NFTAuction{}, err
	}
	return NFTAuction{
		ID:         fmt.Sprintf("%d", i_id),
		Name:       onChainAuction.NFTContract.Hex(),
		StartTime:  "",
		EndTime:    onChainAuction.EndTime.String(),
		LowestBid:  onChainAuction.StartPriceUsd.String(),
		HighestBid: onChainAuction.HighestBidUsd.String(),
		Status:     fmt.Sprintf("%t", onChainAuction.Active),
	}, nil
}

/**
 * 创建拍卖
 */
func (nft *NFTAuctionService) CreateNFTAuction(auction *CreateNFTAuctionRequest) (string, error) {
	if auction == nil {
		return "", fmt.Errorf("创建拍卖参数不能为空")
	}
	if !common.IsHexAddress(auction.NFTContract) {
		return "", fmt.Errorf("NFT 合约地址格式不正确")
	}
	if !isZeroAddressLiteral(auction.PaymentToken) && !common.IsHexAddress(auction.PaymentToken) {
		return "", fmt.Errorf("支付代币地址格式不正确")
	}

	tokenID, ok := new(big.Int).SetString(auction.TokenID, 10)
	if !ok {
		return "", fmt.Errorf("token_id 不是有效的十进制整数")
	}
	startPriceUSD, ok := new(big.Int).SetString(auction.StartPriceUSD, 10)
	if !ok {
		return "", fmt.Errorf("start_price_usd 不是有效的十进制整数")
	}
	durationHours, ok := new(big.Int).SetString(auction.DurationHours, 10)
	if !ok {
		return "", fmt.Errorf("duration_hours 不是有效的十进制整数")
	}

	return block.CreateNFTAuction(&block.ContractAuctionCreate{
		NFTContract:   common.HexToAddress(auction.NFTContract),
		TokenId:       tokenID,
		PaymentToken:  common.HexToAddress(auction.PaymentToken),
		StartPriceUsd: startPriceUSD,
		DurationHours: durationHours,
	})
}

/**
 * 铸造NFT
 */
func (nft *NFTAuctionService) MintNFT(uri string) (string, error) {
	if strings.TrimSpace(uri) == "" {
		return "", fmt.Errorf("uri 不能为空")
	}
	return block.MintNFT(uri)
}

/**
 * 授权NFT
 */
func (nft *NFTAuctionService) ApproveNFT(tokenID string) (string, error) {
	if strings.TrimSpace(tokenID) == "" {
		return "", fmt.Errorf("token_id 不能为空")
	}
	token, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return "", fmt.Errorf("token_id 不是有效的十进制整数")
	}
	return block.ApproveNFT(token)
}

/**
 * 拍卖出价
 */
func (nft *NFTAuctionService) BidNFTAuction(auctionID string, bidAmount string) (string, error) {
	if strings.TrimSpace(auctionID) == "" {
		return "", fmt.Errorf("auction_id 不能为空")
	}
	if strings.TrimSpace(bidAmount) == "" {
		return "", fmt.Errorf("bid_amount 不能为空")
	}

	auction, ok := new(big.Int).SetString(auctionID, 10)
	if !ok {
		return "", fmt.Errorf("auction_id 不是有效的十进制整数")
	}
	amount, ok := new(big.Int).SetString(bidAmount, 10)
	if !ok {
		return "", fmt.Errorf("bid_amount 不是有效的十进制整数")
	}

	return block.BidNFTAuction(auction, amount)
}

func (nft *NFTAuctionService) GetBidLogs(id string) ([]BidLog, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("auction_id 不能为空")
	}
	auctionID, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return nil, fmt.Errorf("auction_id 不是有效的十进制整数")
	}

	onChainLogs, err := block.BidLogs(auctionID)
	if err != nil {
		return nil, err
	}

	logs := make([]BidLog, 0, len(onChainLogs))
	for _, item := range onChainLogs {
		logs = append(logs, BidLog{
			AuctionID:   item.AuctionID.String(),
			Bidder:      item.Bidder.Hex(),
			Amount:      item.Amount.String(),
			AmountUsd:   item.AmountUsd.String(),
			TxHash:      item.TxHash.Hex(),
			BlockNumber: fmt.Sprintf("%d", item.BlockNumber),
			LogIndex:    fmt.Sprintf("%d", item.LogIndex),
		})
	}

	return logs, nil
}

func (nft *NFTAuctionService) GetAuctionLogs(filter *AuctionLogFilter) ([]AuctionLog, error) {
	onChainLogs, err := block.AuctionLogs()
	if err != nil {
		return nil, err
	}

	logs := make([]AuctionLog, 0, len(onChainLogs))
	for _, item := range onChainLogs {
		log := AuctionLog{
			EventType:   item.EventType,
			TxHash:      item.TxHash.Hex(),
			BlockNumber: fmt.Sprintf("%d", item.BlockNumber),
			LogIndex:    fmt.Sprintf("%d", item.LogIndex),
		}
		if item.AuctionID != nil {
			log.AuctionID = item.AuctionID.String()
		}
		if item.Seller != (common.Address{}) {
			log.Seller = item.Seller.Hex()
		}
		if item.NFTContract != (common.Address{}) {
			log.NFTContract = item.NFTContract.Hex()
		}
		if item.TokenId != nil {
			log.TokenID = item.TokenId.String()
		}
		if item.PaymentToken != (common.Address{}) {
			log.PaymentToken = item.PaymentToken.Hex()
		}
		if item.StartPriceUsd != nil {
			log.StartPriceUsd = item.StartPriceUsd.String()
		}
		if item.EndTime != nil {
			log.EndTime = item.EndTime.String()
		}
		if item.Bidder != (common.Address{}) {
			log.Bidder = item.Bidder.Hex()
		}
		if item.Winner != (common.Address{}) {
			log.Winner = item.Winner.Hex()
		}
		if item.Amount != nil {
			log.Amount = item.Amount.String()
		}
		if item.AmountUsd != nil {
			log.AmountUsd = item.AmountUsd.String()
		}
		logs = append(logs, log)
	}

	if filter == nil {
		return logs, nil
	}

	filtered := make([]AuctionLog, 0, len(logs))
	for _, item := range logs {
		if !matchesAuctionLogFilter(item, filter) {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered, nil
}

func matchesAuctionLogFilter(log AuctionLog, filter *AuctionLogFilter) bool {
	if filter == nil {
		return true
	}
	if strings.TrimSpace(filter.EventType) != "" && !strings.EqualFold(strings.TrimSpace(filter.EventType), log.EventType) {
		return false
	}
	if strings.TrimSpace(filter.AuctionID) != "" && log.AuctionID != strings.TrimSpace(filter.AuctionID) {
		return false
	}
	if strings.TrimSpace(filter.Seller) != "" && !strings.EqualFold(strings.TrimSpace(filter.Seller), log.Seller) {
		return false
	}
	if strings.TrimSpace(filter.NFTContract) != "" && !strings.EqualFold(strings.TrimSpace(filter.NFTContract), log.NFTContract) {
		return false
	}
	if strings.TrimSpace(filter.TokenID) != "" && log.TokenID != strings.TrimSpace(filter.TokenID) {
		return false
	}
	if strings.TrimSpace(filter.Bidder) != "" && !strings.EqualFold(strings.TrimSpace(filter.Bidder), log.Bidder) {
		return false
	}
	if strings.TrimSpace(filter.Winner) != "" && !strings.EqualFold(strings.TrimSpace(filter.Winner), log.Winner) {
		return false
	}
	if strings.TrimSpace(filter.TxHash) != "" && !strings.EqualFold(strings.TrimSpace(filter.TxHash), log.TxHash) {
		return false
	}
	if strings.TrimSpace(filter.FromBlock) != "" {
		fromBlock, err := strconv.ParseUint(strings.TrimSpace(filter.FromBlock), 10, 64)
		if err != nil || parseUint64(log.BlockNumber) < fromBlock {
			return false
		}
	}
	if strings.TrimSpace(filter.ToBlock) != "" {
		toBlock, err := strconv.ParseUint(strings.TrimSpace(filter.ToBlock), 10, 64)
		if err != nil || parseUint64(log.BlockNumber) > toBlock {
			return false
		}
	}
	return true
}

func parseUint64(value string) uint64 {
	parsed, _ := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	return parsed
}

func isZeroAddressLiteral(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "0x0" || normalized == "0x" {
		return true
	}
	if !strings.HasPrefix(normalized, "0x") {
		return false
	}
	return len(normalized) == 42 && normalized[2:] == "0000000000000000000000000000000000000000"
}
func (nft *NFTAuctionService) GetTotalNFTAuctions() (OverviewAution, error) {
	return OverviewAution{}, nil
}
func (nft *NFTAuctionService) EndNFTAuction(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("auction_id 不能为空")
	}
	auction, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return fmt.Errorf("auction_id 不是有效的十进制整数")
	}
	return block.EndNFTAuction(auction)
}
func NewNFTAuctionService() NFTAuctionServiceInterface {
	return &NFTAuctionService{}
}
