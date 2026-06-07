package router

import (
	"NFTAuctionServer/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(g *gin.Engine) {
	// 定义路由
	g.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	g.GET("/nftauction", api.GetNFTAuctionList)
	g.GET("/nftauction/:id", api.GetNFTAuctionByID)
	g.GET("/nftauction/:id/bids", api.GetBidLogs)
	g.GET("/nftauction/logs", api.GetAuctionLogs)
	g.POST("/nftauction", api.CreateNFTAuction)
	g.POST("/nftauction/end", api.EndNFTAuction)
	g.POST("/nftauction/bid", api.BidNFTAuction)
	g.POST("/nft/mint", api.MintNFT)
	g.POST("/nft/approve", api.ApproveNFT)
	g.GET("/swagger.json", api.GetSwaggerDoc)
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))
}
