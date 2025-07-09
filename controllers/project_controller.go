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

// SetDefaultProject は指定されたプロジェクトをデフォルトに設定し、他のプロジェクトのデフォルト設定を解除します
func SetDefaultProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id") // デフォルトに設定するプロジェクトのID

	// トランザクションを開始
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 既存のすべてのプロジェクトのIsDefaultをfalseに設定
	if err := tx.Model(&models.Project{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset default projects"})
		return
	}

	// 指定されたプロジェクトのIsDefaultをtrueに設定
	var project models.Project
	if err := tx.First(&project, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	project.IsDefault = true
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default project"})
		return
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/projects")
}
