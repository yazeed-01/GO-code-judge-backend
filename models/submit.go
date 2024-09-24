package models

type Submission struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	ProblemID uint   `json:"problem_id"`
	UserID    uint   `json:"user_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	Token     string `json:"token"`
}
