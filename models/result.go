package models

import (
	"time"
)

type Result struct {
	ID                uint      `gorm:"primaryKey"` // This will auto-increment the ID
	UserID            uint      `json:"user_id"`
	ContestID         uint      `json:"contest_id"`
	ProblemID         uint      `json:"problem_id"`
	LanguageID        int       `json:"language_id"`
	SourceCode        string    `json:"source_code"`
	Input             string    `json:"input"`
	Output            string    `json:"output"`
	ErrorOutput       string    `json:"error_output"`
	StatusID          uint      `json:"status_id"`
	StatusDescription string    `json:"status_desc"`
	MemoryUsed        int       `json:"memory_used"`
	CreatedAt         time.Time `json:"created_at"`
	FinishedAt        time.Time `json:"finished_at"`
	ResultMessage     string    `json:"result_message"`
}
