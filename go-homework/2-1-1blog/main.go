// @title           Blog API
// @version         1.0
// @description     Blog 后台接口文档
// @host      localhost:8080
// @BasePath  /
package main

import (
	"blog/config"
	"blog/models"
	"blog/router"

	_ "blog/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	g := gin.Default()
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config.InitConfig()
	models.InitDB()
	models.AutoMigrate()
	router.InitRouter(g)
	g.Run(":8080")
}
