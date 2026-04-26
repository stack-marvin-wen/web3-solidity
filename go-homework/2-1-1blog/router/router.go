package router

import (
	"blog/view"

	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	// 用户相关路由
	userGroup := g.Group("/users")
	{
		userGroup.POST("/", view.CreateUser)
		userGroup.GET("querybyemail", view.GetUserByEmail)
		userGroup.GET("/:id", view.GetUserByID)
		userGroup.PUT("/", view.UpdateUser)
		userGroup.DELETE("/:id", view.DeleteUser)
		// 其他用户相关路由，如 GET /users/:id, PUT /users/:id, DELETE /users/:id 等
	}
	articleGroup := g.Group("/articles")
	{
		articleGroup.POST("/", view.CreateArticle)
		articleGroup.GET("/:id", view.GetArticleByID)
		articleGroup.PUT("/:id", view.UpdateArticle)
		articleGroup.DELETE("/:id", view.DeleteArticle)
		// 其他文章相关路由，如 GET /articles/:id, PUT /articles/:id, DELETE /articles/:id 等
	}
	tagGroup := g.Group("/tags")
	{
		tagGroup.POST("/", view.CreateTag)
		tagGroup.GET("/:id", view.GetTagByID)
		tagGroup.PUT("/:id", view.UpdateTag)
		tagGroup.DELETE("/:id", view.DeleteTag)
		// 其他标签相关路由，如 GET /tags/:id, PUT /tags/:id, DELETE /tags/:id 等
	}
}
