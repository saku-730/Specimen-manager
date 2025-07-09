package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"unique;not null"` // ユーザー名 (重複不可、必須)
	IsDefault bool   `gorm:"default:false"`   // デフォルトユーザーかどうか
}
