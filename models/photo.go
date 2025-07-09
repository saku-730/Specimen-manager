package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	SpecimenID uint
	Path       string `gorm:"not null"`
}
