// @Summary 文章相关接口
// @Tags 文章接口
// @Description 文章相关接口，包括创建、查询、更新和删除文章
package view

import (
	"blog/models"
	"blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArticleRequest struct {
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	UserID  int      `json:"user_id" binding:"required"`
}
type UpdateArticleRequest struct {
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	UserID  int      `json:"user_id" binding:"required"`
}

var asvc = services.NewArticleService()

// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param data body view.CreateArticleRequest true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /articles [post]
func CreateArticle(ctx *gin.Context) {
	var req CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article := &models.Article{
		Content: req.Content,
		Title:   req.Title,
		UserID:  uint(req.UserID),
	}
	msg, err := asvc.CreateArticle(article)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "文章创建成功", "article": msg})
}

// @Summary 获取文章
// @Description 根据ID获取文章信息
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path string true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /articles/{id} [get]
func GetArticleByID(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	article, err := asvc.GetArticleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"article": article})
}

// @Summary 更新文章
// @Description 根据ID更新文章信息
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path string true "文章ID"
// @Param article body UpdateArticleRequest true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /articles/{id} [put]
func UpdateArticle(ctx *gin.Context) {
	var req UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ctx.Params.ByName("id")
	article, err := asvc.GetArticleByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	article.Content = req.Content
	article.Title = req.Title
	article.UserID = uint(req.UserID)
	err = asvc.UpdateArticle(&article)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "文章更新成功", "article": article})
}

// @Summary 删除文章
// @Description 根据ID删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path string true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /articles/{id} [delete]
func DeleteArticle(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	err := asvc.DeleteArticle(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
