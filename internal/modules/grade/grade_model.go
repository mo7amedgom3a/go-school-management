package grade

import (
	"gorm.io/gorm"

	"school_management/internal/modules/exam"
	"school_management/internal/modules/student"
)

type Grade struct {
	gorm.Model
	StudentID uint    `gorm:"not null" json:"student_id"`
	ExamID    uint    `gorm:"not null" json:"exam_id"`
	Score     float64 `gorm:"not null" json:"score"`

	// Belongs To relationships
	Student student.Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Exam    exam.Exam       `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
}

// TableName specifies the table name for the Grade model
func (Grade) TableName() string {
	return "grades"
}
