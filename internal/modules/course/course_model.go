package course

import (
	"gorm.io/gorm"

	"school_management/internal/modules/department"
	"school_management/internal/modules/teacher"
)

type Course struct {
	gorm.Model
	Name         string `gorm:"not null;size:100" json:"name"`
	Code         string `gorm:"uniqueIndex;not null;size:20" json:"code"`
	Description  string `gorm:"type:text" json:"description"`
	Credits      int    `gorm:"not null;default:3" json:"credits"`
	DepartmentID uint   `gorm:"not null" json:"department_id"`
	TeacherID    uint   `gorm:"not null" json:"teacher_id"`

	// Belongs To relationships
	Department department.Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Teacher    teacher.Teacher       `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// TableName specifies the table name for the Course model
func (Course) TableName() string {
	return "courses"
}
