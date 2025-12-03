package student_courses

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// StudentCourseRepository defines the interface for student course enrollment data access
type StudentCourseRepository interface {
	Create(enrollment *StudentCourse) error
	GetByID(id uint) (*StudentCourse, error)
	GetByIDWithRelations(id uint) (*StudentCourse, error)
	GetAll(limit, offset int) ([]StudentCourse, error)
	GetByStudent(studentID uint) ([]StudentCourse, error)
	GetByCourse(courseID uint) ([]StudentCourse, error)
	GetByStudentAndCourse(studentID, courseID uint) (*StudentCourse, error)
	GetEnrolledAfter(date time.Time) ([]StudentCourse, error)
	Delete(id uint) error
}

// studentCourseRepository implements StudentCourseRepository
type studentCourseRepository struct {
	db *gorm.DB
}

// NewStudentCourseRepository creates a new student course repository with dependency injection
func NewStudentCourseRepository(db *gorm.DB) StudentCourseRepository {
	return &studentCourseRepository{db: db}
}

// Create creates a new enrollment
func (r *studentCourseRepository) Create(enrollment *StudentCourse) error {
	if err := r.db.Create(enrollment).Error; err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	return nil
}

// GetByID retrieves an enrollment by ID
func (r *studentCourseRepository) GetByID(id uint) (*StudentCourse, error) {
	var enrollment StudentCourse
	if err := r.db.First(&enrollment, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return &enrollment, nil
}

// GetByIDWithRelations retrieves an enrollment with student and course preloaded
func (r *studentCourseRepository) GetByIDWithRelations(id uint) (*StudentCourse, error) {
	var enrollment StudentCourse
	if err := r.db.Preload("Student").Preload("Course").First(&enrollment, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollment with relations: %w", err)
	}
	return &enrollment, nil
}

// GetAll retrieves all enrollments with pagination
func (r *studentCourseRepository) GetAll(limit, offset int) ([]StudentCourse, error) {
	var enrollments []StudentCourse
	if err := r.db.Limit(limit).Offset(offset).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	return enrollments, nil
}

// GetByStudent retrieves all enrollments for a student
func (r *studentCourseRepository) GetByStudent(studentID uint) ([]StudentCourse, error) {
	var enrollments []StudentCourse
	if err := r.db.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollments by student: %w", err)
	}
	return enrollments, nil
}

// GetByCourse retrieves all enrollments for a course
func (r *studentCourseRepository) GetByCourse(courseID uint) ([]StudentCourse, error) {
	var enrollments []StudentCourse
	if err := r.db.Where("course_id = ?", courseID).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollments by course: %w", err)
	}
	return enrollments, nil
}

// GetByStudentAndCourse retrieves a specific enrollment
func (r *studentCourseRepository) GetByStudentAndCourse(studentID, courseID uint) (*StudentCourse, error) {
	var enrollment StudentCourse
	if err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).
		First(&enrollment).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return &enrollment, nil
}

// GetEnrolledAfter retrieves enrollments created after a specific date
func (r *studentCourseRepository) GetEnrolledAfter(date time.Time) ([]StudentCourse, error) {
	var enrollments []StudentCourse
	if err := r.db.Where("enrollment_date > ?", date).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to get enrollments after date: %w", err)
	}
	return enrollments, nil
}

// Delete soft deletes an enrollment
func (r *studentCourseRepository) Delete(id uint) error {
	if err := r.db.Delete(&StudentCourse{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	return nil
}
