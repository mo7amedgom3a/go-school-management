package grade

import (
	"fmt"

	"gorm.io/gorm"
)

// GradeRepository defines the interface for grade data access
type GradeRepository interface {
	Create(grade *Grade) error
	GetByID(id uint) (*Grade, error)
	GetByIDWithRelations(id uint) (*Grade, error)
	GetAll(limit, offset int) ([]Grade, error)
	GetByStudent(studentID uint) ([]Grade, error)
	GetByExam(examID uint) ([]Grade, error)
	GetStudentAverage(studentID uint) (float64, error)
	GetExamAverage(examID uint) (float64, error)
	Update(grade *Grade) error
	Delete(id uint) error
}

// gradeRepository implements GradeRepository
type gradeRepository struct {
	db *gorm.DB
}

// NewGradeRepository creates a new grade repository with dependency injection
func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db: db}
}

// Create creates a new grade
func (r *gradeRepository) Create(grade *Grade) error {
	if err := r.db.Create(grade).Error; err != nil {
		return fmt.Errorf("failed to create grade: %w", err)
	}
	return nil
}

// GetByID retrieves a grade by ID
func (r *gradeRepository) GetByID(id uint) (*Grade, error) {
	var grade Grade
	if err := r.db.First(&grade, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get grade: %w", err)
	}
	return &grade, nil
}

// GetByIDWithRelations retrieves a grade with student and exam preloaded
func (r *gradeRepository) GetByIDWithRelations(id uint) (*Grade, error) {
	var grade Grade
	if err := r.db.Preload("Student").Preload("Exam").First(&grade, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get grade with relations: %w", err)
	}
	return &grade, nil
}

// GetAll retrieves all grades with pagination
func (r *gradeRepository) GetAll(limit, offset int) ([]Grade, error) {
	var grades []Grade
	if err := r.db.Limit(limit).Offset(offset).Find(&grades).Error; err != nil {
		return nil, fmt.Errorf("failed to get grades: %w", err)
	}
	return grades, nil
}

// GetByStudent retrieves all grades for a student
func (r *gradeRepository) GetByStudent(studentID uint) ([]Grade, error) {
	var grades []Grade
	if err := r.db.Where("student_id = ?", studentID).Find(&grades).Error; err != nil {
		return nil, fmt.Errorf("failed to get grades by student: %w", err)
	}
	return grades, nil
}

// GetByExam retrieves all grades for an exam
func (r *gradeRepository) GetByExam(examID uint) ([]Grade, error) {
	var grades []Grade
	if err := r.db.Where("exam_id = ?", examID).Find(&grades).Error; err != nil {
		return nil, fmt.Errorf("failed to get grades by exam: %w", err)
	}
	return grades, nil
}

// GetStudentAverage calculates the average grade for a student
func (r *gradeRepository) GetStudentAverage(studentID uint) (float64, error) {
	var avg float64
	if err := r.db.Model(&Grade{}).
		Where("student_id = ?", studentID).
		Select("AVG(score)").
		Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate student average: %w", err)
	}
	return avg, nil
}

// GetExamAverage calculates the average grade for an exam
func (r *gradeRepository) GetExamAverage(examID uint) (float64, error) {
	var avg float64
	if err := r.db.Model(&Grade{}).
		Where("exam_id = ?", examID).
		Select("AVG(score)").
		Scan(&avg).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate exam average: %w", err)
	}
	return avg, nil
}

// Update updates a grade
func (r *gradeRepository) Update(grade *Grade) error {
	if err := r.db.Save(grade).Error; err != nil {
		return fmt.Errorf("failed to update grade: %w", err)
	}
	return nil
}

// Delete soft deletes a grade
func (r *gradeRepository) Delete(id uint) error {
	if err := r.db.Delete(&Grade{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete grade: %w", err)
	}
	return nil
}
