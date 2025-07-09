package controllers

import (
	"GO-specimen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProjects はプロジェクト一覧を表示します
func GetProjects(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var projects []models.Project
	db.Find(&projects)
	c.HTML(http.StatusOK, "projects_index.html", gin.H{"Projects": projects})
}

// CreateProject は新しいプロジェクトを作成します
func CreateProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name cannot be empty"})
		return
	}

	project := models.Project{Name: name}
	if err := db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/projects")
}

// DeleteProject はプロジェクトを削除します
func DeleteProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// プロジェクトが存在するか確認
	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// プロジェクトを削除
	if err := db.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/projects")
}
