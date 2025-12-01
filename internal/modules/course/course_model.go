package course

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name         string `gorm:"not null;size:100" json:"name"`
	Code         string `gorm:"uniqueIndex;not null;size:20" json:"code"`
	Description  string `gorm:"type:text" json:"description"`
	Credits      int    `gorm:"not null;default:3" json:"credits"`
	DepartmentID uint   `gorm:"not null" json:"department_id"`
	TeacherID    uint   `gorm:"not null" json:"teacher_id"`

	// Relationships
	Department  interface{}   `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Teacher     interface{}   `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Students    []interface{} `gorm:"many2many:student_courses" json:"students,omitempty"`
	Attendances []interface{} `gorm:"foreignKey:CourseID" json:"attendances,omitempty"`
	Homework    []interface{} `gorm:"foreignKey:CourseID" json:"homework,omitempty"`
	Exams       []interface{} `gorm:"foreignKey:CourseID" json:"exams,omitempty"`
}

// TableName specifies the table name for the Course model
func (Course) TableName() string {
	return "courses"
}
