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

	// Relationships (ignored during migration, used with Preload)
	Department     Department      `gorm:"-" json:"department,omitempty"`
	Teacher        Teacher         `gorm:"-" json:"teacher,omitempty"`
	StudentCourses []StudentCourse `gorm:"-" json:"student_courses,omitempty"`
	Attendances    []Attendance    `gorm:"-" json:"attendances,omitempty"`
	Homework       []Homework      `gorm:"-" json:"homework,omitempty"`
	Exams          []Exam          `gorm:"-" json:"exams,omitempty"`
}

// Placeholder types
type Department struct{}
type Teacher struct{}
type StudentCourse struct{}
type Attendance struct{}
type Homework struct{}
type Exam struct{}

// TableName specifies the table name for the Course model
func (Course) TableName() string {
	return "courses"
}
