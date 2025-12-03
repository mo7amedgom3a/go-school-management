package student

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// StudentRepository defines the interface for student data access
type StudentRepository interface {
	Create(student *Student) error
	GetByID(id uint) (*Student, error)
	GetAll(limit, offset int) ([]Student, error)
	GetByEmail(email string) (*Student, error)
	Search(query string, limit int) ([]Student, error)
	GetEnrolledBefore(date time.Time) ([]Student, error)
	Update(student *Student) error
	Delete(id uint) error
}

// studentRepository implements StudentRepository
type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new student repository with dependency injection
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

// Create creates a new student
func (r *studentRepository) Create(student *Student) error {
	if err := r.db.Create(student).Error; err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}
	return nil
}

// GetByID retrieves a student by ID
func (r *studentRepository) GetByID(id uint) (*Student, error) {
	var student Student
	if err := r.db.First(&student, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	return &student, nil
}

// GetAll retrieves all students with pagination
func (r *studentRepository) GetAll(limit, offset int) ([]Student, error) {
	var students []Student
	if err := r.db.Limit(limit).Offset(offset).Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}
	return students, nil
}

// GetByEmail retrieves a student by email
func (r *studentRepository) GetByEmail(email string) (*Student, error) {
	var student Student
	if err := r.db.Where("email = ?", email).First(&student).Error; err != nil {
		return nil, fmt.Errorf("failed to get student by email: %w", err)
	}
	return &student, nil
}

// Search searches students by name or email
func (r *studentRepository) Search(query string, limit int) ([]Student, error) {
	var students []Student
	searchPattern := "%" + query + "%"
	if err := r.db.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
		searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to search students: %w", err)
	}
	return students, nil
}

// GetEnrolledBefore retrieves students enrolled before a specific date
func (r *studentRepository) GetEnrolledBefore(date time.Time) ([]Student, error) {
	var students []Student
	if err := r.db.Where("enrollment_date < ?", date).Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to get students enrolled before date: %w", err)
	}
	return students, nil
}

// Update updates a student
func (r *studentRepository) Update(student *Student) error {
	if err := r.db.Save(student).Error; err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}
	return nil
}

// Delete soft deletes a student
func (r *studentRepository) Delete(id uint) error {
	if err := r.db.Delete(&Student{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	return nil
}
