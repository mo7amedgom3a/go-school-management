package course

import "time"

// CreateCourseRequest represents the request body for creating a course
type CreateCourseRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=100"`
	Code         string `json:"code" binding:"required,min=2,max=20"`
	Description  string `json:"description" binding:"omitempty,max=500"`
	Credits      int    `json:"credits" binding:"required,min=1,max=10"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	TeacherID    uint   `json:"teacher_id" binding:"required"`
}

// UpdateCourseRequest represents the request body for updating a course
type UpdateCourseRequest struct {
	Name         string `json:"name" binding:"omitempty,min=2,max=100"`
	Code         string `json:"code" binding:"omitempty,min=2,max=20"`
	Description  string `json:"description" binding:"omitempty,max=500"`
	Credits      int    `json:"credits" binding:"omitempty,min=1,max=10"`
	DepartmentID uint   `json:"department_id" binding:"omitempty"`
	TeacherID    uint   `json:"teacher_id" binding:"omitempty"`
}

// CourseResponse represents the response body for course data
type CourseResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	Credits      int       `json:"credits"`
	DepartmentID uint      `json:"department_id"`
	TeacherID    uint      `json:"teacher_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
