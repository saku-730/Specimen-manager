package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name string `gorm:"unique;not null"` // プロジェクト名 (重複不可、必須)
}
