package grade

import "gorm.io/gorm"

type Grade struct {
	gorm.Model
	StudentID uint    `gorm:"not null" json:"student_id"`
	ExamID    uint    `gorm:"not null" json:"exam_id"`
	Score     float64 `gorm:"not null" json:"score"`

	// Relationships (ignored during migration, used with Preload)
	Student Student `gorm:"-" json:"student,omitempty"`
	Exam    Exam    `gorm:"-" json:"exam,omitempty"`
}

// Placeholder types
type Student struct{}
type Exam struct{}

// TableName specifies the table name for the Grade model
func (Grade) TableName() string {
	return "grades"
}
