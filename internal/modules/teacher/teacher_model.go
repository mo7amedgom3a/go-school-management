package teacher

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	FirstName    string `gorm:"not null;size:50" json:"first_name"`
	LastName     string `gorm:"not null;size:50" json:"last_name"`
	Email        string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Phone        string `gorm:"size:20" json:"phone"`
	DepartmentID uint   `gorm:"not null" json:"department_id"`

	// Relationships (ignored during migration, used with Preload)
	Department Department `gorm:"-" json:"department,omitempty"`
	Courses    []Course   `gorm:"-" json:"courses,omitempty"`
}

// Placeholder types
type Department struct{}
type Course struct{}

// TableName specifies the table name for the Teacher model
func (Teacher) TableName() string {
	return "teachers"
}
