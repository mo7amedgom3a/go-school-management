package student

import "time"

// CreateStudentRequest represents the request body for creating a student
type CreateStudentRequest struct {
	FirstName      string `json:"first_name" binding:"required,min=2,max=50"`
	LastName       string `json:"last_name" binding:"required,min=2,max=50"`
	Email          string `json:"email" binding:"required,email,max=100"`
	Phone          string `json:"phone" binding:"omitempty,max=20"`
	DateOfBirth    string `json:"date_of_birth" binding:"required"`   // Format: YYYY-MM-DD
	EnrollmentDate string `json:"enrollment_date" binding:"required"` // Format: YYYY-MM-DD
}

// UpdateStudentRequest represents the request body for updating a student
type UpdateStudentRequest struct {
	FirstName      string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName       string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Email          string `json:"email" binding:"omitempty,email,max=100"`
	Phone          string `json:"phone" binding:"omitempty,max=20"`
	DateOfBirth    string `json:"date_of_birth" binding:"omitempty"`   // Format: YYYY-MM-DD
	EnrollmentDate string `json:"enrollment_date" binding:"omitempty"` // Format: YYYY-MM-DD
}

// StudentResponse represents the response body for student data
type StudentResponse struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
