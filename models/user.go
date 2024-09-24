package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UserID   uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	FullName string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}
