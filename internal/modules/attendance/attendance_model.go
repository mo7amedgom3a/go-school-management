package attendance

import (
	"time"

	"gorm.io/gorm"

	"school_management/internal/modules/course"
	"school_management/internal/modules/student"
)

// AttendanceStatus represents the status of attendance
type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "present"
	AttendanceAbsent  AttendanceStatus = "absent"
	AttendanceLate    AttendanceStatus = "late"
)

type Attendance struct {
	gorm.Model
	StudentID uint             `gorm:"not null" json:"student_id"`
	CourseID  uint             `gorm:"not null" json:"course_id"`
	Date      time.Time        `gorm:"type:date;not null" json:"date"`
	Status    AttendanceStatus `gorm:"type:varchar(20);not null;default:'present'" json:"status"`

	// Belongs To relationships
	Student student.Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course  course.Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for the Attendance model
func (Attendance) TableName() string {
	return "attendances"
}
