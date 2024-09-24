package models

import (
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model

	ContestID        uint   `gorm:"primaryKey"`
	ContestName      string `gorm:"unique;not null"`
	ParticipantCount int    `gorm:"default:0;not null"`
}
