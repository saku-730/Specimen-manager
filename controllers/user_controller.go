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

// SetDefaultUser は指定されたユーザーをデフォルトに設定し、他のユーザーのデフォルト設定を解除します
func SetDefaultUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id") // デフォルトに設定するユーザーのID

	// トランザクションを開始
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 既存のすべてのユーザーのIsDefaultをfalseに設定
	if err := tx.Model(&models.User{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset default users"})
		return
	}

	// 指定されたユーザーのIsDefaultをtrueに設定
	var user models.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.IsDefault = true
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default user"})
		return
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Redirect(http.StatusFound, "/settings/users")
}
