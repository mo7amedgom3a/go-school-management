package grade

import "time"

// CreateGradeRequest represents the request body for creating a grade
type CreateGradeRequest struct {
	StudentID uint    `json:"student_id" binding:"required"`
	ExamID    uint    `json:"exam_id" binding:"required"`
	Score     float64 `json:"score" binding:"required,min=0"`
}

// UpdateGradeRequest represents the request body for updating a grade
type UpdateGradeRequest struct {
	Score float64 `json:"score" binding:"required,min=0"`
}

// GradeResponse represents the response body for grade data
type GradeResponse struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"student_id"`
	ExamID    uint      `json:"exam_id"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
