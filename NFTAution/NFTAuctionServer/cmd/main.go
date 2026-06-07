package main

import (
	"NFTAuctionServer/config"
	"NFTAuctionServer/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	config.InitConfig()
	router.InitRouter(g)
	g.Run(fmt.Sprintf(":%v", config.Config.AppConfig.Port))
}
