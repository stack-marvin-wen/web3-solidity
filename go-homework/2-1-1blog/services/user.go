package services

import (
	"blog/models"
	"fmt"
)

type UserServiceInterface interface {
	CreateUser(user *models.User) (string, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}
type UserService struct {
}

func (usvc *UserService) CreateUser(user *models.User) (string, error) {
	if err := models.DB.Create(user).Error; err != nil {
		return fmt.Sprintf("用户创建失败：%v", err), err
	}
	return "创建用户成功！", nil
}
func (usvc *UserService) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("用户未找到：%v", err)
	}
	return user, nil
}
func (usvc *UserService) GetUserByID(id string) (models.User, error) {
	var user models.User
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("用户未找到：%v", err)
	}
	return user, nil
}
func (usvc *UserService) UpdateUser(user *models.User) error {

	if err := models.DB.Save(user).Error; err != nil {
		return fmt.Errorf("用户更新失败：%v", err)
	}
	return nil
}
func (usvc *UserService) DeleteUser(id string) error {
	if err := models.DB.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("用户删除失败：%v", err)
	}
	models.DB.Unscoped().Delete(&models.Article{}, "user_id = ?", id)
	return nil
}
func NewUserService() UserServiceInterface {
	return &UserService{}
}
