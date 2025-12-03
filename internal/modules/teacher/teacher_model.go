package teacher

import (
	"gorm.io/gorm"

	"school_management/internal/modules/department"
)

type Teacher struct {
	gorm.Model
	FirstName    string `gorm:"not null;size:50" json:"first_name"`
	LastName     string `gorm:"not null;size:50" json:"last_name"`
	Email        string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Phone        string `gorm:"size:20" json:"phone"`
	DepartmentID uint   `gorm:"not null" json:"department_id"`

	// Belongs To relationship
	Department department.Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// TableName specifies the table name for the Teacher model
func (Teacher) TableName() string {
	return "teachers"
}
