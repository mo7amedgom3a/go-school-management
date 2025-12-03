package teacher

import (
	"fmt"

	"gorm.io/gorm"
)

// TeacherRepository defines the interface for teacher data access
type TeacherRepository interface {
	Create(teacher *Teacher) error
	GetByID(id uint) (*Teacher, error)
	GetByIDWithDepartment(id uint) (*Teacher, error)
	GetAll(limit, offset int) ([]Teacher, error)
	GetByDepartment(deptID uint) ([]Teacher, error)
	GetByEmail(email string) (*Teacher, error)
	Update(teacher *Teacher) error
	Delete(id uint) error
}

// teacherRepository implements TeacherRepository
type teacherRepository struct {
	db *gorm.DB
}

// NewTeacherRepository creates a new teacher repository with dependency injection
func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

// Create creates a new teacher
func (r *teacherRepository) Create(teacher *Teacher) error {
	if err := r.db.Create(teacher).Error; err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}
	return nil
}

// GetByID retrieves a teacher by ID
func (r *teacherRepository) GetByID(id uint) (*Teacher, error) {
	var teacher Teacher
	if err := r.db.First(&teacher, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}
	return &teacher, nil
}

// GetByIDWithDepartment retrieves a teacher by ID with department preloaded
func (r *teacherRepository) GetByIDWithDepartment(id uint) (*Teacher, error) {
	var teacher Teacher
	if err := r.db.Preload("Department").First(&teacher, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get teacher with department: %w", err)
	}
	return &teacher, nil
}

// GetAll retrieves all teachers with pagination
func (r *teacherRepository) GetAll(limit, offset int) ([]Teacher, error) {
	var teachers []Teacher
	if err := r.db.Limit(limit).Offset(offset).Find(&teachers).Error; err != nil {
		return nil, fmt.Errorf("failed to get teachers: %w", err)
	}
	return teachers, nil
}

// GetByDepartment retrieves all teachers in a department
func (r *teacherRepository) GetByDepartment(deptID uint) ([]Teacher, error) {
	var teachers []Teacher
	if err := r.db.Where("department_id = ?", deptID).Find(&teachers).Error; err != nil {
		return nil, fmt.Errorf("failed to get teachers by department: %w", err)
	}
	return teachers, nil
}

// GetByEmail retrieves a teacher by email
func (r *teacherRepository) GetByEmail(email string) (*Teacher, error) {
	var teacher Teacher
	if err := r.db.Where("email = ?", email).First(&teacher).Error; err != nil {
		return nil, fmt.Errorf("failed to get teacher by email: %w", err)
	}
	return &teacher, nil
}

// Update updates a teacher
func (r *teacherRepository) Update(teacher *Teacher) error {
	if err := r.db.Save(teacher).Error; err != nil {
		return fmt.Errorf("failed to update teacher: %w", err)
	}
	return nil
}

// Delete soft deletes a teacher
func (r *teacherRepository) Delete(id uint) error {
	if err := r.db.Delete(&Teacher{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}
	return nil
}
