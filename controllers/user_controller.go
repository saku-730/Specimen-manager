package controllers

import (
	"GO-specimen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUsers はユーザー一覧を表示します
func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)
	c.HTML(http.StatusOK, "users_index.html", gin.H{"Users": users})
}

// CreateUser は新しいユーザーを作成します
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	name := c.PostForm("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name cannot be empty"})
		return
	}

	user := models.User{Name: name}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/users")
}

// DeleteUser はユーザーを削除します
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// ユーザーが存在するか確認
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// ユーザーを削除
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/users")
}
