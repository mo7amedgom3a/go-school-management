package students_homework

import (
	"time"

	"gorm.io/gorm"
)

// HomeworkStatus represents the status of homework submission
type HomeworkStatus string

const (
	HomeworkPending   HomeworkStatus = "pending"
	HomeworkSubmitted HomeworkStatus = "submitted"
	HomeworkGraded    HomeworkStatus = "graded"
)

type StudentHomework struct {
	gorm.Model
	StudentID      uint           `gorm:"not null;uniqueIndex:idx_student_homework" json:"student_id"`
	HomeworkID     uint           `gorm:"not null;uniqueIndex:idx_student_homework" json:"homework_id"`
	SubmissionDate *time.Time     `gorm:"type:timestamp" json:"submission_date"`
	Score          *float64       `json:"score"`
	Status         HomeworkStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// Relationships
	Student  interface{} `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Homework interface{} `gorm:"foreignKey:HomeworkID" json:"homework,omitempty"`
}

// TableName specifies the table name for the StudentHomework model
func (StudentHomework) TableName() string {
	return "students_homework"
}
