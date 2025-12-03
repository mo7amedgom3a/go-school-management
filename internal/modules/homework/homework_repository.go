package homework

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// HomeworkRepository defines the interface for homework data access
type HomeworkRepository interface {
	Create(homework *Homework) error
	GetByID(id uint) (*Homework, error)
	GetByIDWithCourse(id uint) (*Homework, error)
	GetAll(limit, offset int) ([]Homework, error)
	GetByCourse(courseID uint) ([]Homework, error)
	GetUpcoming(limit int) ([]Homework, error)
	GetOverdue() ([]Homework, error)
	Update(homework *Homework) error
	Delete(id uint) error
}

// homeworkRepository implements HomeworkRepository
type homeworkRepository struct {
	db *gorm.DB
}

// NewHomeworkRepository creates a new homework repository with dependency injection
func NewHomeworkRepository(db *gorm.DB) HomeworkRepository {
	return &homeworkRepository{db: db}
}

// Create creates a new homework assignment
func (r *homeworkRepository) Create(homework *Homework) error {
	if err := r.db.Create(homework).Error; err != nil {
		return fmt.Errorf("failed to create homework: %w", err)
	}
	return nil
}

// GetByID retrieves a homework assignment by ID
func (r *homeworkRepository) GetByID(id uint) (*Homework, error) {
	var homework Homework
	if err := r.db.First(&homework, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get homework: %w", err)
	}
	return &homework, nil
}

// GetByIDWithCourse retrieves a homework assignment with course preloaded
func (r *homeworkRepository) GetByIDWithCourse(id uint) (*Homework, error) {
	var homework Homework
	if err := r.db.Preload("Course").First(&homework, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get homework with course: %w", err)
	}
	return &homework, nil
}

// GetAll retrieves all homework assignments with pagination
func (r *homeworkRepository) GetAll(limit, offset int) ([]Homework, error) {
	var homeworks []Homework
	if err := r.db.Limit(limit).Offset(offset).Find(&homeworks).Error; err != nil {
		return nil, fmt.Errorf("failed to get homeworks: %w", err)
	}
	return homeworks, nil
}

// GetByCourse retrieves all homework assignments for a course
func (r *homeworkRepository) GetByCourse(courseID uint) ([]Homework, error) {
	var homeworks []Homework
	if err := r.db.Where("course_id = ?", courseID).Find(&homeworks).Error; err != nil {
		return nil, fmt.Errorf("failed to get homeworks by course: %w", err)
	}
	return homeworks, nil
}

// GetUpcoming retrieves upcoming homework assignments (due in the future)
func (r *homeworkRepository) GetUpcoming(limit int) ([]Homework, error) {
	var homeworks []Homework
	if err := r.db.Where("due_date > ?", time.Now()).
		Order("due_date ASC").
		Limit(limit).
		Find(&homeworks).Error; err != nil {
		return nil, fmt.Errorf("failed to get upcoming homeworks: %w", err)
	}
	return homeworks, nil
}

// GetOverdue retrieves overdue homework assignments
func (r *homeworkRepository) GetOverdue() ([]Homework, error) {
	var homeworks []Homework
	if err := r.db.Where("due_date < ?", time.Now()).
		Order("due_date DESC").
		Find(&homeworks).Error; err != nil {
		return nil, fmt.Errorf("failed to get overdue homeworks: %w", err)
	}
	return homeworks, nil
}

// Update updates a homework assignment
func (r *homeworkRepository) Update(homework *Homework) error {
	if err := r.db.Save(homework).Error; err != nil {
		return fmt.Errorf("failed to update homework: %w", err)
	}
	return nil
}

// Delete soft deletes a homework assignment
func (r *homeworkRepository) Delete(id uint) error {
	if err := r.db.Delete(&Homework{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete homework: %w", err)
	}
	return nil
}
