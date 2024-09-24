package models

import (
	"gorm.io/gorm"
)

type Language struct {
	gorm.Model

	LanguageID   uint   `gorm:"primaryKey"`
	LanguageName string `gorm:"unique;not null"`
	Function     string `gorm:"not null"`
}
