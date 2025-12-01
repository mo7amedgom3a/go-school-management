package teacher

import "time"

// CreateTeacherRequest represents the request body for creating a teacher
type CreateTeacherRequest struct {
	FirstName    string `json:"first_name" binding:"required,min=2,max=50"`
	LastName     string `json:"last_name" binding:"required,min=2,max=50"`
	Email        string `json:"email" binding:"required,email,max=100"`
	Phone        string `json:"phone" binding:"omitempty,max=20"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

// UpdateTeacherRequest represents the request body for updating a teacher
type UpdateTeacherRequest struct {
	FirstName    string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName     string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Email        string `json:"email" binding:"omitempty,email,max=100"`
	Phone        string `json:"phone" binding:"omitempty,max=20"`
	DepartmentID uint   `json:"department_id" binding:"omitempty"`
}

// TeacherResponse represents the response body for teacher data
type TeacherResponse struct {
	ID           uint      `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	DepartmentID uint      `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
