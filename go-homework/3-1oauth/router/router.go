package router

import (
	"oauth/middleware"
	"oauth/views"

	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	userGroup := g.Group("/user")
	{
		userGroup.POST("/login", views.Login)
	}
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.POST("/register", views.CreateUser)
	}
}
