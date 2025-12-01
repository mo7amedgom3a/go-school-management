package department

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Name        string `gorm:"not null;size:100" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Relationships (ignored during migration, used with Preload)
	Teachers []Teacher `gorm:"-" json:"teachers,omitempty"`
	Courses  []Course  `gorm:"-" json:"courses,omitempty"`
}

// Placeholder types to avoid circular imports - will be replaced when querying
type Teacher struct{}
type Course struct{}

// TableName specifies the table name for the Department model
func (Department) TableName() string {
	return "departments"
}
