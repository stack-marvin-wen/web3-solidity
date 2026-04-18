package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func uploadSingleFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filename := file.Filename
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "文件上传成功", "filename": filename})
}
func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join("uploads", filename)
	if _, err := os.Stat(filepath); err != nil {
		c.JSON(400, gin.H{"error": "文件不存在"})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}
func uploadMutipleFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	var filenames []string
	for _, file := range files {
		filename := file.Filename
		dst := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		filenames = append(filenames, filename)
	}
	c.JSON(200, gin.H{"message": "文件上传成功", "filenames": filenames})
}
func main() {
	r := gin.Default()
	// 	# 创建测试文件
	// echo "Hello, World!" > test.txt

	// # 上传文件
	// curl -X POST http://localhost:8080/uploadsinglefile  -F "file=@test.txt"
	r.POST("/uploadsinglefile", uploadSingleFile)
	// 	# 创建多个测试文件
	// echo "File 1" > file1.txt
	// echo "File 2" > file2.txt
	// echo "File 3" > file3.txt
	// # 上传多个文件
	// curl -X POST http://localhost:8080/uploadmutiplefile  -F "files=@file1.txt" -F "files=@file2.txt" -F "files=@file3.txt"
	r.POST("/uploadmutiplefile", uploadMutipleFile)
	r.GET("/download/:filename", downloadFile)
	r.Static("/uploads", "./uploads")
	r.StaticFS("/files", http.Dir("./uploads"))
	r.Run(":8080")
}
