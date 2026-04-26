package service

import (
	"fmt"
	"oauth/models"
	"oauth/utils"
)

type UserServiceInterface interface {
	CreateUser(user *models.User) (string, error)
	GetUserByID(id string) (models.User, error)
	Login(email string, password string) (string, error)
}

type UserService struct {
}

func NewUserService() UserServiceInterface {
	return &UserService{}
}

func (usvc *UserService) CreateUser(user *models.User) (string, error) {
	err := models.DB.Create(user).Error
	if err != nil {
		return "增加失败", err
	}
	return "增加成功", nil
}

func (usvc *UserService) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := models.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (usvc *UserService) Login(email, password string) (string, error) {
	var user models.User
	err := models.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return "账号不存在", err
	}
	if user.UserPWD != password {
		return "密码错误", fmt.Errorf("密码错误")
	}
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "生成token失败", err
	}
	return token, nil
}
