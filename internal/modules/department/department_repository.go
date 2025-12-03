package department

import (
	"fmt"

	"gorm.io/gorm"
)

// DepartmentRepository defines the interface for department data access
type DepartmentRepository interface {
	Create(dept *Department) error
	GetByID(id uint) (*Department, error)
	GetAll(limit, offset int) ([]Department, error)
	Update(dept *Department) error
	Delete(id uint) error
	Search(name string) ([]Department, error)
}

// departmentRepository implements DepartmentRepository
type departmentRepository struct {
	db *gorm.DB
}
// Department Constructor
// NewDepartmentRepository creates a new department repository with dependency injection
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

// Create creates a new department
func (r *departmentRepository) Create(dept *Department) error {
	if err := r.db.Create(dept).Error; err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}
	return nil
}

// GetByID retrieves a department by ID
func (r *departmentRepository) GetByID(id uint) (*Department, error) {
	var dept Department
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return &dept, nil
}

// GetAll retrieves all departments with pagination
func (r *departmentRepository) GetAll(limit, offset int) ([]Department, error) {
	var departments []Department
	query := r.db.Limit(limit).Offset(offset)
	if err := query.Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}
	return departments, nil
}

// Update updates a department
func (r *departmentRepository) Update(dept *Department) error {
	if err := r.db.Save(dept).Error; err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}
	return nil
}

// Delete soft deletes a department
func (r *departmentRepository) Delete(id uint) error {
	if err := r.db.Delete(&Department{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}
	return nil
}

// Search searches departments by name
func (r *departmentRepository) Search(name string) ([]Department, error) {
	var departments []Department
	if err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to search departments: %w", err)
	}
	return departments, nil
}
