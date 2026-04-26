package middleware

import (
	"net/http"
	"oauth/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取token
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "请求未携带token",
			})
			ctx.Abort()
			return
		}
		// 解析token
		_, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "token解析失败",
			})
			ctx.Abort()
			return
		}
	}
}
