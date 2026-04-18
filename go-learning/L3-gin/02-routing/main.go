package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 基础路由
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "hello"})
	})

	// 路由携带参数
	// 单个参数
	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, gin.H{"user_id": id})
	})
	// 多个参数
	r.GET("/user/:id/profile/:name", func(ctx *gin.Context) {
		id := ctx.Param("id")
		name := ctx.Param("name")
		ctx.JSON(http.StatusOK, gin.H{"user_id": id, "profile_name": name})
	})
	// 通配符参数
	r.GET("/files/*filepath", func(ctx *gin.Context) {
		filepath := ctx.Param("filepath")
		ctx.JSON(http.StatusOK, gin.H{
			"filepath": filepath,
		})
	})
	// 查询参数
	r.GET("/search", func(ctx *gin.Context) {
		keyword := ctx.Query("kw")
		page := ctx.DefaultQuery("page", "1")
		size := ctx.DefaultQuery("size", "10") //参数默认值
		ctx.JSON(http.StatusOK, gin.H{
			"keyword": keyword,
			"page":    page,
			"size":    size,
		})
	})
	// 表单参数
	/*
		测试：curl -X POST http://localhost:8080/login -H "Content-Type: application/x-www-form-urlencoded" -d "username=admin&userpwd=123456&remember=true"
	*/
	r.POST("/login", func(ctx *gin.Context) {
		uname := ctx.PostForm("username")
		upwd := ctx.PostForm("userpwd")
		remember := ctx.DefaultPostForm("remember", "false")
		ctx.JSON(http.StatusOK, gin.H{
			"username": uname,
			"userpwd":  upwd,
			"remember": remember,
		})
	})

	// JSON表单参数
	/*
		测试：curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json"  -d '{"name":"Alice","email":"alice@example.com","age":25}'
	*/
	r.POST("/api/users", func(ctx *gin.Context) {
		type CreateUserReq struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required"`
			Age   int    `json:"age" binding:"gte=0,lte=130"`
		}
		var req CreateUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "创建用户成功",
			"id":      1,
			"name":    req.Name,
			"email":   req.Email,
			"age":     req.Age,
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
		v1.GET("/user/:id", getUserByID)
	}
}
func getUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"users": []gin.H{
			{
				"id":   1,
				"name": "张三",
				"age":  18,
			},
			{
				"id":   2,
				"name": "李四",
				"age":  19,
			},
		},
	})
}
func getUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取用户详情成功",
		"user": gin.H{
			"id":   id,
			"name": "张三",
			"age":  18,
		},
	})
}
