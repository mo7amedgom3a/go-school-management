package attendance

import (
	"time"

	"gorm.io/gorm"
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

	// Relationships (ignored during migration, used with Preload)
	Student Student `gorm:"-" json:"student,omitempty"`
	Course  Course  `gorm:"-" json:"course,omitempty"`
}

// Placeholder types
type Student struct{}
type Course struct{}

// TableName specifies the table name for the Attendance model
func (Attendance) TableName() string {
	return "attendances"
}
