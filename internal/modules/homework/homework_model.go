package homework

import (
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	Title       string    `gorm:"not null;size:200" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	DueDate     time.Time `gorm:"type:timestamp;not null" json:"due_date"`
	MaxScore    float64   `gorm:"not null;default:100" json:"max_score"`

	// Relationships
	Course      interface{}   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Submissions []interface{} `gorm:"foreignKey:HomeworkID" json:"submissions,omitempty"`
}

// TableName specifies the table name for the Homework model
func (Homework) TableName() string {
	return "homework"
}
