package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func createroute() *gin.Engine {
	return gin.Default()
}

func middlewareRegister(r *gin.Engine) {
	r.Use(loggerMiddleware())
}
func routeRegister(r *gin.Engine) {
	r.GET("/ping", ping)
}
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Printf("[%s] %s %d %v\n", c.Request.Method, path, status, latency)
	}
}
func main() {
	r := createroute()
	middlewareRegister(r)
	routeRegister(r)
	r.Run(":8080")
}
func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ping")
}
