package exam

import (
	"time"

	"gorm.io/gorm"

	"school_management/internal/modules/course"
)

type Exam struct {
	gorm.Model
	Title    string    `gorm:"not null;size:200" json:"title"`
	CourseID uint      `gorm:"not null" json:"course_id"`
	ExamDate time.Time `gorm:"type:timestamp;not null" json:"exam_date"`
	Duration int       `gorm:"not null;comment:Duration in minutes" json:"duration"`
	MaxScore float64   `gorm:"not null;default:100" json:"max_score"`

	// Belongs To relationship
	Course course.Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for the Exam model
func (Exam) TableName() string {
	return "exams"
}
