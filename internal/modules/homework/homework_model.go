package homework

import (
	"time"

	"gorm.io/gorm"

	"school_management/internal/modules/course"
)

type Homework struct {
	gorm.Model
	Title       string    `gorm:"not null;size:200" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	DueDate     time.Time `gorm:"type:timestamp;not null" json:"due_date"`
	MaxScore    float64   `gorm:"not null;default:100" json:"max_score"`

	// Belongs To relationship
	Course course.Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for the Homework model
func (Homework) TableName() string {
	return "homework"
}
