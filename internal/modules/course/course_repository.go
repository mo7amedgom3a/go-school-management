package course

import (
	"fmt"

	"gorm.io/gorm"
)

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	Create(course *Course) error
	GetByID(id uint) (*Course, error)
	GetByIDWithRelations(id uint) (*Course, error)
	GetAll(limit, offset int) ([]Course, error)
	GetByDepartment(deptID uint) ([]Course, error)
	GetByTeacher(teacherID uint) ([]Course, error)
	GetByCode(code string) (*Course, error)
	Update(course *Course) error
	Delete(id uint) error
}

// courseRepository implements CourseRepository
type courseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository with dependency injection
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

// Create creates a new course
func (r *courseRepository) Create(course *Course) error {
	if err := r.db.Create(course).Error; err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

// GetByID retrieves a course by ID
func (r *courseRepository) GetByID(id uint) (*Course, error) {
	var course Course
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	return &course, nil
}

// GetByIDWithRelations retrieves a course by ID with department and teacher preloaded
func (r *courseRepository) GetByIDWithRelations(id uint) (*Course, error) {
	var course Course
	if err := r.db.Preload("Department").Preload("Teacher").First(&course, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get course with relations: %w", err)
	}
	return &course, nil
}

// GetAll retrieves all courses with pagination
func (r *courseRepository) GetAll(limit, offset int) ([]Course, error) {
	var courses []Course
	if err := r.db.Limit(limit).Offset(offset).Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

// GetByDepartment retrieves all courses in a department
func (r *courseRepository) GetByDepartment(deptID uint) ([]Course, error) {
	var courses []Course
	if err := r.db.Where("department_id = ?", deptID).Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("failed to get courses by department: %w", err)
	}
	return courses, nil
}

// GetByTeacher retrieves all courses taught by a teacher
func (r *courseRepository) GetByTeacher(teacherID uint) ([]Course, error) {
	var courses []Course
	if err := r.db.Where("teacher_id = ?", teacherID).Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("failed to get courses by teacher: %w", err)
	}
	return courses, nil
}

// GetByCode retrieves a course by code
func (r *courseRepository) GetByCode(code string) (*Course, error) {
	var course Course
	if err := r.db.Where("code = ?", code).First(&course).Error; err != nil {
		return nil, fmt.Errorf("failed to get course by code: %w", err)
	}
	return &course, nil
}

// Update updates a course
func (r *courseRepository) Update(course *Course) error {
	if err := r.db.Save(course).Error; err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}

// Delete soft deletes a course
func (r *courseRepository) Delete(id uint) error {
	if err := r.db.Delete(&Course{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	return nil
}
