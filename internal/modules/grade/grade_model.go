package grade

import "gorm.io/gorm"

type Grade struct {
	gorm.Model
	StudentID uint    `gorm:"not null" json:"student_id"`
	ExamID    uint    `gorm:"not null" json:"exam_id"`
	Score     float64 `gorm:"not null" json:"score"`

	// Relationships
	Student interface{} `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Exam    interface{} `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
}

// TableName specifies the table name for the Grade model
func (Grade) TableName() string {
	return "grades"
}
