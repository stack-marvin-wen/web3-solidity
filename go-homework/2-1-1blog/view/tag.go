package view

import (
	"blog/models"
	"blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tagService = services.NewTagService()

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}
type UpdateTagRequest struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// @Summary 标签相关接口
// @Description 标签相关接口，包括创建、查询、更新和删除标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param data body view.CreateTagRequest true "标签信息"
// @Success 200 {object} map[string]interface{}
// @Router /tags [post]
func CreateTag(ctx *gin.Context) {
	var tagRequest CreateTagRequest
	if err := ctx.ShouldBindJSON(&tagRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag := &models.Tag{
		Name: tagRequest.Name,
	}
	msg, err := tagService.CreateTag(tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": msg})
}

// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Param data body view.UpdateTagRequest true "标签信息"
// @Success 200 {object} map[string]interface{}
// @Router /tags [put]
func UpdateTag(ctx *gin.Context) {
	var tagRequest UpdateTagRequest
	if err := ctx.ShouldBindJSON(&tagRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := tagService.GetTagByID(tagRequest.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	tag.Name = tagRequest.Name
	err = tagService.UpdateTag(&tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "标签更新成功", "tag": tag})

}

// @Summary 删除标签
// @Description 根据ID删除标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {object} map[string]interface{}
// @Router /tags/{id} [delete]
func DeleteTag(ctx *gin.Context) {
	id := ctx.Param("id")
	err := tagService.DeleteTag(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "标签删除成功"})
}

// @Summary 获取标签
// @Description 根据ID获取标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {object} map[string]interface{}
// @Router /tags/{id} [get]
func GetTagByID(ctx *gin.Context) {
	id := ctx.Param("id")
	tag, err := tagService.GetTagByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"tag": tag})
}
