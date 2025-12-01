package department

import "time"

// CreateDepartmentRequest represents the request body for creating a department
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// UpdateDepartmentRequest represents the request body for updating a department
type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// DepartmentResponse represents the response body for department data
type DepartmentResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
