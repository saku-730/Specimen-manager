package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name      string `gorm:"unique;not null"` // プロジェクト名 (重複不可、必須)
	IsDefault bool   `gorm:"default:false"`   // デフォルトプロジェクトかどうか
}
