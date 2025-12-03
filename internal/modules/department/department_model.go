package department

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Name        string `gorm:"not null;size:100" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

// TableName specifies the table name for the Department model
func (Department) TableName() string {
	return "departments"
}
