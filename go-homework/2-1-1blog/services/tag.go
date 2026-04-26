package services

import "blog/models"

type TagServiceInterface interface {
	CreateTag(tag *models.Tag) (string, error)
	GetTagByID(id string) (models.Tag, error)
	UpdateTag(tag *models.Tag) error
	DeleteTag(id string) error
}
type TagService struct {
}

func (tsvc *TagService) CreateTag(tag *models.Tag) (string, error) {
	if err := models.DB.Create(tag).Error; err != nil {
		return "标签创建失败！", err
	}
	return "创建标签成功！", nil
}
func (tsvc *TagService) GetTagByID(id string) (models.Tag, error) {
	var tag models.Tag
	if err := models.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return models.Tag{}, err
	}
	return tag, nil
}
func (tsvc *TagService) UpdateTag(tag *models.Tag) error {
	if err := models.DB.Save(tag).Error; err != nil {
		return err
	}
	return nil
}
func (tsvc *TagService) DeleteTag(id string) error {
	if err := models.DB.Delete(&models.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}
func NewTagService() TagServiceInterface {
	return &TagService{}
}
