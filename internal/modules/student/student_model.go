package student

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName      string    `gorm:"not null;size:50" json:"first_name"`
	LastName       string    `gorm:"not null;size:50" json:"last_name"`
	Email          string    `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Phone          string    `gorm:"size:20" json:"phone"`
	DateOfBirth    time.Time `gorm:"type:date" json:"date_of_birth"`
	EnrollmentDate time.Time `gorm:"type:date;not null" json:"enrollment_date"`

	// Relationships (ignored during migration, used with Preload)
	StudentCourses  []StudentCourse   `gorm:"-" json:"student_courses,omitempty"`
	StudentHomework []StudentHomework `gorm:"-" json:"student_homework,omitempty"`
	Attendances     []Attendance      `gorm:"-" json:"attendances,omitempty"`
	Grades          []Grade           `gorm:"-" json:"grades,omitempty"`
}

// Placeholder types
type StudentCourse struct{}
type StudentHomework struct{}
type Attendance struct{}
type Grade struct{}

// TableName specifies the table name for the Student model
func (Student) TableName() string {
	return "students"
}
