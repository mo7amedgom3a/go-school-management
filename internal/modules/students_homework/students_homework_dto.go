package students_homework

import "time"

// SubmitHomeworkRequest represents the request body for submitting homework
type SubmitHomeworkRequest struct {
	StudentID      uint   `json:"student_id" binding:"required"`
	HomeworkID     uint   `json:"homework_id" binding:"required"`
	SubmissionDate string `json:"submission_date" binding:"required"` // Format: YYYY-MM-DD HH:MM:SS
}

// GradeHomeworkRequest represents the request body for grading homework
type GradeHomeworkRequest struct {
	StudentID  uint    `json:"student_id" binding:"required"`
	HomeworkID uint    `json:"homework_id" binding:"required"`
	Score      float64 `json:"score" binding:"required,min=0"`
}

// StudentHomeworkResponse represents the response body for student homework data
type StudentHomeworkResponse struct {
	ID             uint      `json:"id"`
	StudentID      uint      `json:"student_id"`
	HomeworkID     uint      `json:"homework_id"`
	SubmissionDate *string   `json:"submission_date"`
	Score          *float64  `json:"score"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
