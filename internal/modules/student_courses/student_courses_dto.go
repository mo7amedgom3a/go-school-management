package student_courses

import "time"

// EnrollStudentRequest represents the request body for enrolling a student in a course
type EnrollStudentRequest struct {
	StudentID      uint   `json:"student_id" binding:"required"`
	CourseID       uint   `json:"course_id" binding:"required"`
	EnrollmentDate string `json:"enrollment_date" binding:"required"` // Format: YYYY-MM-DD
}

// StudentCourseResponse represents the response body for student course enrollment data
type StudentCourseResponse struct {
	ID             uint      `json:"id"`
	StudentID      uint      `json:"student_id"`
	CourseID       uint      `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
