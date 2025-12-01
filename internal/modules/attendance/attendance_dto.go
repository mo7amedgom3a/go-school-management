package attendance

import "time"

// CreateAttendanceRequest represents the request body for creating an attendance record
type CreateAttendanceRequest struct {
	StudentID uint   `json:"student_id" binding:"required"`
	CourseID  uint   `json:"course_id" binding:"required"`
	Date      string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Status    string `json:"status" binding:"required,oneof=present absent late"`
}

// UpdateAttendanceRequest represents the request body for updating an attendance record
type UpdateAttendanceRequest struct {
	Status string `json:"status" binding:"omitempty,oneof=present absent late"`
}

// AttendanceResponse represents the response body for attendance data
type AttendanceResponse struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"student_id"`
	CourseID  uint      `json:"course_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
