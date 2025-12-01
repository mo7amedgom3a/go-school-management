package student_courses

import (
	"time"

	"gorm.io/gorm"
)

type StudentCourse struct {
	gorm.Model
	StudentID      uint      `gorm:"not null;uniqueIndex:idx_student_course" json:"student_id"`
	CourseID       uint      `gorm:"not null;uniqueIndex:idx_student_course" json:"course_id"`
	EnrollmentDate time.Time `gorm:"type:date;not null" json:"enrollment_date"`

	// Relationships (ignored during migration, used with Preload)
	Student Student `gorm:"-" json:"student,omitempty"`
	Course  Course  `gorm:"-" json:"course,omitempty"`
}

// Placeholder types
type Student struct{}
type Course struct{}

// TableName specifies the table name for the StudentCourse model
func (StudentCourse) TableName() string {
	return "student_courses"
}
