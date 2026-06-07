package api

import (
	"fmt"
	"net/http"

	"NFTAuctionServer/service"

	"github.com/gin-gonic/gin"
)

var autionService = service.NewNFTAuctionService()

type NFTAuctionListResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []service.NFTAuction `json:"data"`
}

type NFTAuctionDetailResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    service.NFTAuction `json:"data"`
}

type CreateNFTAuctionResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    CreateNFTAuctionResult `json:"data"`
}

type CreateNFTAuctionResult struct {
	TxHash string `json:"tx_hash"`
}

type MintNFTResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    CreateNFTAuctionResult `json:"data"`
}

type ApproveNFTResponse = MintNFTResponse

type BidNFTAuctionResponse = MintNFTResponse

type EndNFTAuctionResponse = MintNFTResponse

type BidLogsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []service.BidLog `json:"data"`
}

type AuctionLogsResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []service.AuctionLog `json:"data"`
}

// GetNFTAuctionList godoc
// @Summary 获取拍卖列表
// @Description 从链上读取当前所有 NFT 拍卖
// @Tags NFTAuction
// @Produce json
// @Success 200 {object} NFTAuctionListResponse
// @Failure 500 {object} NFTAuctionListResponse
// @Router /nftauction [get]
func GetNFTAuctionList(c *gin.Context) {
	data, err := autionService.GetNFTAuctionList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取拍卖列表失败, %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取拍卖列表成功",
		"data":    data,
	})
}

// GetNFTAuctionByID godoc
// @Summary 获取拍卖详情
// @Description 根据拍卖 ID 查询链上拍卖信息
// @Tags NFTAuction
// @Produce json
// @Param id path string true "拍卖 ID"
// @Success 200 {object} NFTAuctionDetailResponse
// @Failure 500 {object} NFTAuctionDetailResponse
// @Router /nftauction/{id} [get]
func GetNFTAuctionByID(c *gin.Context) {
	id := c.Param("id")
	data, err := autionService.GetNFTAuctionByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取拍卖详情失败, %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取拍卖详情成功",
		"data":    data,
	})
}

// CreateNFTAuction godoc
// @Summary 创建拍卖
// @Description 调用链上合约创建一个新的 NFT 拍卖
// @Tags NFTAuction
// @Accept json
// @Produce json
// @Param auction body service.CreateNFTAuctionRequest true "创建拍卖请求"
// @Success 200 {object} CreateNFTAuctionResponse
// @Failure 400 {object} CreateNFTAuctionResponse
// @Failure 500 {object} CreateNFTAuctionResponse
// @Router /nftauction [post]
func CreateNFTAuction(c *gin.Context) {
	var req service.CreateNFTAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误, %v", err),
		})
		return
	}

	txHash, err := autionService.CreateNFTAuction(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("创建拍卖失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建拍卖成功",
		"data": CreateNFTAuctionResult{
			TxHash: txHash,
		},
	})
}

// MintNFT godoc
// @Summary 铸造 NFT
// @Description 使用当前配置的钱包在 NFT 合约上铸造一个 NFT
// @Tags NFTInstance
// @Accept json
// @Produce json
// @Param body body service.MintNFTRequest true "铸造请求"
// @Success 200 {object} MintNFTResponse
// @Failure 400 {object} MintNFTResponse
// @Failure 500 {object} MintNFTResponse
// @Router /nft/mint [post]
func MintNFT(c *gin.Context) {
	var req service.MintNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误, %v", err),
		})
		return
	}

	txHash, err := autionService.MintNFT(req.URI)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("铸造 NFT 失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "铸造 NFT 成功",
		"data": CreateNFTAuctionResult{
			TxHash: txHash,
		},
	})
}

// ApproveNFT godoc
// @Summary 授权 NFT 给拍卖合约
// @Description 将指定 tokenId 授权给当前配置的拍卖合约
// @Tags NFTInstance
// @Accept json
// @Produce json
// @Param body body service.ApproveNFTRequest true "授权请求"
// @Success 200 {object} ApproveNFTResponse
// @Failure 400 {object} ApproveNFTResponse
// @Failure 500 {object} ApproveNFTResponse
// @Router /nft/approve [post]
func ApproveNFT(c *gin.Context) {
	var req service.ApproveNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误, %v", err),
		})
		return
	}

	txHash, err := autionService.ApproveNFT(req.TokenID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("授权 NFT 失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "授权 NFT 成功",
		"data": CreateNFTAuctionResult{
			TxHash: txHash,
		},
	})
}

// BidNFTAuction godoc
// @Summary 拍卖出价
// @Description 根据拍卖的 paymentToken 自动选择 ETH 出价或 ERC20 出价
// @Tags NFTAuction
// @Accept json
// @Produce json
// @Param body body service.BidNFTAuctionRequest true "出价请求"
// @Success 200 {object} BidNFTAuctionResponse
// @Failure 400 {object} BidNFTAuctionResponse
// @Failure 500 {object} BidNFTAuctionResponse
// @Router /nftauction/bid [post]
func BidNFTAuction(c *gin.Context) {
	var req service.BidNFTAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误, %v", err),
		})
		return
	}

	txHash, err := autionService.BidNFTAuction(req.AuctionID, req.BidAmount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("拍卖出价失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "拍卖出价成功",
		"data": CreateNFTAuctionResult{
			TxHash: txHash,
		},
	})
}

// GetBidLogs godoc
// @Summary 获取出价日志
// @Description 获取指定拍卖的链上出价事件日志
// @Tags NFTAuction
// @Produce json
// @Param id path string true "拍卖 ID"
// @Success 200 {object} BidLogsResponse
// @Failure 500 {object} BidLogsResponse
// @Router /nftauction/{id}/bids [get]
func GetBidLogs(c *gin.Context) {
	id := c.Param("id")
	data, err := autionService.GetBidLogs(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取出价日志失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取出价日志成功",
		"data":    data,
	})
}

// GetAuctionLogs godoc
// @Summary 获取拍卖全量日志
// @Description 获取所有创建拍卖、出价和结束拍卖事件日志
// @Tags NFTAuction
// @Produce json
// @Param event_type query string false "事件类型：AuctionCreated, HighestBidIncreased, AuctionEnded"
// @Param auction_id query string false "拍卖 ID"
// @Param seller query string false "卖家地址"
// @Param nft_contract query string false "NFT 合约地址"
// @Param token_id query string false "NFT Token ID"
// @Param bidder query string false "出价人地址"
// @Param winner query string false "成交人地址"
// @Param tx_hash query string false "交易哈希"
// @Param from_block query string false "起始区块号"
// @Param to_block query string false "结束区块号"
// @Success 200 {object} AuctionLogsResponse
// @Failure 500 {object} AuctionLogsResponse
// @Router /nftauction/logs [get]
func GetAuctionLogs(c *gin.Context) {
	filter := &service.AuctionLogFilter{
		EventType:   c.Query("event_type"),
		AuctionID:   c.Query("auction_id"),
		Seller:      c.Query("seller"),
		NFTContract: c.Query("nft_contract"),
		TokenID:     c.Query("token_id"),
		Bidder:      c.Query("bidder"),
		Winner:      c.Query("winner"),
		TxHash:      c.Query("tx_hash"),
		FromBlock:   c.Query("from_block"),
		ToBlock:     c.Query("to_block"),
	}
	data, err := autionService.GetAuctionLogs(filter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取拍卖日志失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取拍卖日志成功",
		"data":    data,
	})
}

// EndNFTAuction godoc
// @Summary 结束拍卖
// @Description 结束指定拍卖并结算链上状态
// @Tags NFTAuction
// @Accept json
// @Produce json
// @Param body body service.EndNFTAuctionRequest true "结束拍卖请求"
// @Success 200 {object} EndNFTAuctionResponse
// @Failure 400 {object} EndNFTAuctionResponse
// @Failure 500 {object} EndNFTAuctionResponse
// @Router /nftauction/end [post]
func EndNFTAuction(c *gin.Context) {
	var req service.EndNFTAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("请求参数错误, %v", err),
		})
		return
	}

	if err := autionService.EndNFTAuction(req.AuctionID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("结束拍卖失败, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "结束拍卖成功",
	})
}
