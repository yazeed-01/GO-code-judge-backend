package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model

	ProblemID          uint   `gorm:"primaryKey"`
	ProblemName        string `gorm:"unique;not null"`
	ProblemDescription string `gorm:"not null"`
	Defficulty         string `gorm:"not null"`
	Tags               string `gorm:"not null"`
	ContestID          uint   `gorm:"not null"`
	InputExample       string `gorm:"not null"`
	OutputExample      string `gorm:"not null"`
	TestCaseInput      string `gorm:"not null"`
	TestCaseOutput     string `gorm:"not null"`
}
