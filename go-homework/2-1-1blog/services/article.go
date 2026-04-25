package services

import "blog/models"

type ArticleServiceInterface interface {
	CreateArticle(article *models.Article) (string, error)
	GetArticleByID(id string) (models.Article, error)
	UpdateArticle(article *models.Article) error
	DeleteArticle(id string) error
}
type ArticleService struct {
}

func (asvc *ArticleService) CreateArticle(article *models.Article) (string, error) {
	if err := models.DB.Create(article).Error; err != nil {
		return "文章创建失败！", err
	}
	return "创建文章成功！", nil
}
func (asvc *ArticleService) GetArticleByID(id string) (models.Article, error) {
	var article models.Article
	if err := models.DB.First(&article, id).Error; err != nil {
		return models.Article{}, err
	}
	return article, nil
}
func (asvc *ArticleService) UpdateArticle(article *models.Article) error {
	if err := models.DB.Save(article).Error; err != nil {
		return err
	}
	return nil
}
func (asvc *ArticleService) DeleteArticle(id string) error {
	if err := models.DB.Delete(&models.Article{}, id).Error; err != nil {
		return err
	}
	return nil
}

func NewArticleService() ArticleServiceInterface {
	return &ArticleService{}
}
