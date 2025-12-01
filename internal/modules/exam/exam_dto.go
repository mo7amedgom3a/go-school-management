package exam

import "time"

// CreateExamRequest represents the request body for creating an exam
type CreateExamRequest struct {
	Title    string  `json:"title" binding:"required,min=2,max=200"`
	CourseID uint    `json:"course_id" binding:"required"`
	ExamDate string  `json:"exam_date" binding:"required"`               // Format: YYYY-MM-DD HH:MM:SS
	Duration int     `json:"duration" binding:"required,min=15,max=300"` // Minutes
	MaxScore float64 `json:"max_score" binding:"required,min=1,max=1000"`
}

// UpdateExamRequest represents the request body for updating an exam
type UpdateExamRequest struct {
	Title    string  `json:"title" binding:"omitempty,min=2,max=200"`
	ExamDate string  `json:"exam_date" binding:"omitempty"`               // Format: YYYY-MM-DD HH:MM:SS
	Duration int     `json:"duration" binding:"omitempty,min=15,max=300"` // Minutes
	MaxScore float64 `json:"max_score" binding:"omitempty,min=1,max=1000"`
}

// ExamResponse represents the response body for exam data
type ExamResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CourseID  uint      `json:"course_id"`
	ExamDate  time.Time `json:"exam_date"`
	Duration  int       `json:"duration"`
	MaxScore  float64   `json:"max_score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
