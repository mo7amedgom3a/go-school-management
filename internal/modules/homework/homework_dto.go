package homework

import "time"

// CreateHomeworkRequest represents the request body for creating homework
type CreateHomeworkRequest struct {
	Title       string  `json:"title" binding:"required,min=2,max=200"`
	Description string  `json:"description" binding:"omitempty,max=1000"`
	CourseID    uint    `json:"course_id" binding:"required"`
	DueDate     string  `json:"due_date" binding:"required"` // Format: YYYY-MM-DD HH:MM:SS
	MaxScore    float64 `json:"max_score" binding:"required,min=1,max=1000"`
}

// UpdateHomeworkRequest represents the request body for updating homework
type UpdateHomeworkRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=2,max=200"`
	Description string  `json:"description" binding:"omitempty,max=1000"`
	DueDate     string  `json:"due_date" binding:"omitempty"` // Format: YYYY-MM-DD HH:MM:SS
	MaxScore    float64 `json:"max_score" binding:"omitempty,min=1,max=1000"`
}

// HomeworkResponse represents the response body for homework data
type HomeworkResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CourseID    uint      `json:"course_id"`
	DueDate     time.Time `json:"due_date"`
	MaxScore    float64   `json:"max_score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
