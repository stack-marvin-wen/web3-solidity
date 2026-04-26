// @Summary 用户相关接口
// @Description 用户相关接口，包括创建、查询、更新和删除用户
// @Tags 用户接口
package view

import (
	"blog/models"
	"blog/services"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type UpdateUserRequest struct {
	Id       int    `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

var userService = services.NewUserService()

// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body view.CreateUserRequest true "用户信息"
// @Success 200 {object} map[string]interface{}
// @Router /users [post]
func CreateUser(ctx *gin.Context) {
	request := CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password, // 注意：实际项目中应对密码进行哈希处理
	}
	msg, err := userService.CreateUser(user)
	if err != nil {
		ctx.JSON(500, gin.H{"message": msg})
		return
	}
	ctx.JSON(200, gin.H{"message": msg})
}

// @Summary 根据邮箱获取用户
// @Description 根据邮箱查询用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱"
// @Success 200 {object} map[string]interface{}
// @Router /users/querybyemail [get]
func GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(500, gin.H{"user": models.User{}, "error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"user": user, "error": nil})
}

// @Summary 根据ID获取用户
//
// @Description 根据用户ID查询用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [get]
func GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"user": models.User{}, "error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"user": user, "error": nil})
}

// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body view.UpdateUserRequest true "用户信息"
// @Success 200 {object} map[string]interface{}
// @Router /users [put]
func UpdateUser(ctx *gin.Context) {
	request := UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}
	var user models.User
	if err := models.DB.First(&user, request.Id).Error; err != nil {
		ctx.JSON(404, gin.H{"message": "用户未找到"})
		return
	}
	user.Name = request.Name
	user.Password = request.Password
	user.Email = request.Email
	err := userService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "用户更新成功！"})
}

// @Summary 删除用户
// @Description 删除用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [delete]
func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := userService.DeleteUser(id)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "用户删除成功！"})
}
