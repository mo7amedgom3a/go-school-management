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
}

// TableName specifies the table name for the Student model
func (Student) TableName() string {
	return "students"
}
