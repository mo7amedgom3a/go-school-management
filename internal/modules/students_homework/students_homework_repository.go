package students_homework

import (
	"fmt"

	"gorm.io/gorm"
)

// StudentHomeworkRepository defines the interface for student homework submission data access
type StudentHomeworkRepository interface {
	Create(submission *StudentHomework) error
	GetByID(id uint) (*StudentHomework, error)
	GetByIDWithRelations(id uint) (*StudentHomework, error)
	GetAll(limit, offset int) ([]StudentHomework, error)
	GetByStudent(studentID uint) ([]StudentHomework, error)
	GetByHomework(homeworkID uint) ([]StudentHomework, error)
	GetByStudentAndHomework(studentID, homeworkID uint) (*StudentHomework, error)
	GetByStatus(status HomeworkStatus) ([]StudentHomework, error)
	GetPendingByStudent(studentID uint) ([]StudentHomework, error)
	Update(submission *StudentHomework) error
	Delete(id uint) error
}

// studentHomeworkRepository implements StudentHomeworkRepository
type studentHomeworkRepository struct {
	db *gorm.DB
}

// NewStudentHomeworkRepository creates a new student homework repository with dependency injection
func NewStudentHomeworkRepository(db *gorm.DB) StudentHomeworkRepository {
	return &studentHomeworkRepository{db: db}
}

// Create creates a new homework submission
func (r *studentHomeworkRepository) Create(submission *StudentHomework) error {
	if err := r.db.Create(submission).Error; err != nil {
		return fmt.Errorf("failed to create submission: %w", err)
	}
	return nil
}

// GetByID retrieves a submission by ID
func (r *studentHomeworkRepository) GetByID(id uint) (*StudentHomework, error) {
	var submission StudentHomework
	if err := r.db.First(&submission, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}
	return &submission, nil
}

// GetByIDWithRelations retrieves a submission with student and homework preloaded
func (r *studentHomeworkRepository) GetByIDWithRelations(id uint) (*StudentHomework, error) {
	var submission StudentHomework
	if err := r.db.Preload("Student").Preload("Homework").First(&submission, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get submission with relations: %w", err)
	}
	return &submission, nil
}

// GetAll retrieves all submissions with pagination
func (r *studentHomeworkRepository) GetAll(limit, offset int) ([]StudentHomework, error) {
	var submissions []StudentHomework
	if err := r.db.Limit(limit).Offset(offset).Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get submissions: %w", err)
	}
	return submissions, nil
}

// GetByStudent retrieves all submissions for a student
func (r *studentHomeworkRepository) GetByStudent(studentID uint) ([]StudentHomework, error) {
	var submissions []StudentHomework
	if err := r.db.Where("student_id = ?", studentID).Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get submissions by student: %w", err)
	}
	return submissions, nil
}

// GetByHomework retrieves all submissions for a homework assignment
func (r *studentHomeworkRepository) GetByHomework(homeworkID uint) ([]StudentHomework, error) {
	var submissions []StudentHomework
	if err := r.db.Where("homework_id = ?", homeworkID).Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get submissions by homework: %w", err)
	}
	return submissions, nil
}

// GetByStudentAndHomework retrieves a specific submission
func (r *studentHomeworkRepository) GetByStudentAndHomework(studentID, homeworkID uint) (*StudentHomework, error) {
	var submission StudentHomework
	if err := r.db.Where("student_id = ? AND homework_id = ?", studentID, homeworkID).
		First(&submission).Error; err != nil {
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}
	return &submission, nil
}

// GetByStatus retrieves submissions by status
func (r *studentHomeworkRepository) GetByStatus(status HomeworkStatus) ([]StudentHomework, error) {
	var submissions []StudentHomework
	if err := r.db.Where("status = ?", status).Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get submissions by status: %w", err)
	}
	return submissions, nil
}

// GetPendingByStudent retrieves pending submissions for a student
func (r *studentHomeworkRepository) GetPendingByStudent(studentID uint) ([]StudentHomework, error) {
	var submissions []StudentHomework
	if err := r.db.Where("student_id = ? AND status = ?", studentID, HomeworkPending).
		Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending submissions: %w", err)
	}
	return submissions, nil
}

// Update updates a homework submission
func (r *studentHomeworkRepository) Update(submission *StudentHomework) error {
	if err := r.db.Save(submission).Error; err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}
	return nil
}

// Delete soft deletes a homework submission
func (r *studentHomeworkRepository) Delete(id uint) error {
	if err := r.db.Delete(&StudentHomework{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete submission: %w", err)
	}
	return nil
}
