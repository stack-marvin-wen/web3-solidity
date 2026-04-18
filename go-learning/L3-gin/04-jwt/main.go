package main

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginReqForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key")

func generateJWT(userId int, userName string) (string, error) {
	claims := JWTClaims{
		UserID:   uint(userId),
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}
func parseJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "未授权"})
			c.Abort()
			return
		}
		if auth[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "无效的授权头"})
			c.Abort()
			return
		}
		tokenStr := auth[7:]
		claims, err := parseJWT(tokenStr)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}

}
func main() {
	r := gin.Default()
	r.POST("/login", login)
}

func login(c *gin.Context) {
	var loginReq LoginReqForm
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := generateJWT(1, loginReq.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "生成token失败"})
		return
	}
	c.JSON(200, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}
