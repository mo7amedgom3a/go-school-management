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

	// Relationships
	Student interface{} `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course  interface{} `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for the StudentCourse model
func (StudentCourse) TableName() string {
	return "student_courses"
}
