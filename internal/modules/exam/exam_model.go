package exam

import (
	"time"

	"gorm.io/gorm"
)

type Exam struct {
	gorm.Model
	Title    string    `gorm:"not null;size:200" json:"title"`
	CourseID uint      `gorm:"not null" json:"course_id"`
	ExamDate time.Time `gorm:"type:timestamp;not null" json:"exam_date"`
	Duration int       `gorm:"not null;comment:Duration in minutes" json:"duration"`
	MaxScore float64   `gorm:"not null;default:100" json:"max_score"`

	// Relationships (ignored during migration, used with Preload)
	Course Course  `gorm:"-" json:"course,omitempty"`
	Grades []Grade `gorm:"-" json:"grades,omitempty"`
}

// Placeholder types
type Course struct{}
type Grade struct{}

// TableName specifies the table name for the Exam model
func (Exam) TableName() string {
	return "exams"
}
