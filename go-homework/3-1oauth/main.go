// @title           OAuth Demo API
// @version         1.0
// @description     带 OAuth2 验证的 API 示例
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:8080/user/login
// @scope.read Grants read access
// @scope.write Grants write access
package main

import (
	"oauth/config"
	_ "oauth/docs"
	"oauth/models"
	"oauth/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	g := gin.Default()
	router.InitRouter(g)
	config.InitConfig()
	models.InitDB()
	models.AutoMigrate()
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g.Run(":8080")
}
