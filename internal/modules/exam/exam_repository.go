package exam

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExamRepository defines the interface for exam data access
type ExamRepository interface {
	Create(exam *Exam) error
	GetByID(id uint) (*Exam, error)
	GetByIDWithCourse(id uint) (*Exam, error)
	GetAll(limit, offset int) ([]Exam, error)
	GetByCourse(courseID uint) ([]Exam, error)
	GetUpcoming(limit int) ([]Exam, error)
	GetByDateRange(start, end time.Time) ([]Exam, error)
	Update(exam *Exam) error
	Delete(id uint) error
}

// examRepository implements ExamRepository
type examRepository struct {
	db *gorm.DB
}

// NewExamRepository creates a new exam repository with dependency injection
func NewExamRepository(db *gorm.DB) ExamRepository {
	return &examRepository{db: db}
}

// Create creates a new exam
func (r *examRepository) Create(exam *Exam) error {
	if err := r.db.Create(exam).Error; err != nil {
		return fmt.Errorf("failed to create exam: %w", err)
	}
	return nil
}

// GetByID retrieves an exam by ID
func (r *examRepository) GetByID(id uint) (*Exam, error) {
	var exam Exam
	if err := r.db.First(&exam, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get exam: %w", err)
	}
	return &exam, nil
}

// GetByIDWithCourse retrieves an exam with course preloaded
func (r *examRepository) GetByIDWithCourse(id uint) (*Exam, error) {
	var exam Exam
	if err := r.db.Preload("Course").First(&exam, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get exam with course: %w", err)
	}
	return &exam, nil
}

// GetAll retrieves all exams with pagination
func (r *examRepository) GetAll(limit, offset int) ([]Exam, error) {
	var exams []Exam
	if err := r.db.Limit(limit).Offset(offset).Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to get exams: %w", err)
	}
	return exams, nil
}

// GetByCourse retrieves all exams for a course
func (r *examRepository) GetByCourse(courseID uint) ([]Exam, error) {
	var exams []Exam
	if err := r.db.Where("course_id = ?", courseID).Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to get exams by course: %w", err)
	}
	return exams, nil
}

// GetUpcoming retrieves upcoming exams (scheduled in the future)
func (r *examRepository) GetUpcoming(limit int) ([]Exam, error) {
	var exams []Exam
	if err := r.db.Where("exam_date > ?", time.Now()).
		Order("exam_date ASC").
		Limit(limit).
		Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to get upcoming exams: %w", err)
	}
	return exams, nil
}

// GetByDateRange retrieves exams within a date range
func (r *examRepository) GetByDateRange(start, end time.Time) ([]Exam, error) {
	var exams []Exam
	if err := r.db.Where("exam_date BETWEEN ? AND ?", start, end).Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to get exams by date range: %w", err)
	}
	return exams, nil
}

// Update updates an exam
func (r *examRepository) Update(exam *Exam) error {
	if err := r.db.Save(exam).Error; err != nil {
		return fmt.Errorf("failed to update exam: %w", err)
	}
	return nil
}

// Delete soft deletes an exam
func (r *examRepository) Delete(id uint) error {
	if err := r.db.Delete(&Exam{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete exam: %w", err)
	}
	return nil
}
