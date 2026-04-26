package views

import (
	"oauth/models"
	"oauth/service"

	"github.com/gin-gonic/gin"
)

var userService = service.NewUserService()

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary: CreateUser
// @Description: CreateUser
// @Security OAuth2Password
// @Tags: User
// @Accept: application/json
// @Produce: application/json
// @Param: user body CreateUserRequest true "CreateUserRequest"
// @Success: 200 {object} map[string]interface{}
// @Failure: 400 {object} map[string]interface{}
// @Router: /user/register [post]
func CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := &models.User{
		Email:    req.Email,
		UserPWD:  req.Password,
		Username: req.Username,
	}
	msg, err := userService.CreateUser(user)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": msg, "error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "用户创建成功！", "error": ""})
}

// @Summary: Login
// @Description: Login
// @Tags: User
// @Accept: application/json
// @Produce: application/json
// @Param: user body LoginRequest true "LoginRequest"
// @Success: 200 {object} map[string]interface{}
// @Failure: 400 {object} map[string]interface{}
// @Router: /user/login [post]
func Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	msg, err := userService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"msg": msg, "error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"msg": msg, "error": ""})
}
