package homework

import (
	"fmt"
	"strings"
	"time"
)

// HomeworkService defines the business logic interface
type HomeworkService interface {
	Create(req *CreateHomeworkRequest) (*HomeworkResponse, error)
	GetByID(id uint) (*HomeworkResponse, error)
	GetByCourse(courseID uint) ([]HomeworkResponse, error)
	GetUpcoming(limit int) ([]HomeworkResponse, error)
	Update(id uint, req *UpdateHomeworkRequest) (*HomeworkResponse, error)
	Delete(id uint) error
}

// homeworkService implements HomeworkService
type homeworkService struct {
	repo HomeworkRepository
}

// NewHomeworkService creates a new homework service with DI
func NewHomeworkService(repo HomeworkRepository) HomeworkService {
	return &homeworkService{repo: repo}
}

// Create creates a new homework assignment
func (s *homeworkService) Create(req *CreateHomeworkRequest) (*HomeworkResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse due date
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due date format (use RFC3339): %w", err)
	}

	// Map DTO to Model
	hw := &Homework{
		Title:       req.Title,
		Description: req.Description,
		CourseID:    req.CourseID,
		DueDate:     dueDate,
		MaxScore:    req.MaxScore,
	}

	// Create via repository
	if err := s.repo.Create(hw); err != nil {
		return nil, fmt.Errorf("failed to create homework: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(hw), nil
}

// GetByID retrieves a homework assignment by ID
func (s *homeworkService) GetByID(id uint) (*HomeworkResponse, error) {
	hw, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("homework not found: %w", err)
	}
	return s.toResponseDTO(hw), nil
}

// GetByCourse retrieves homework assignments for a course
func (s *homeworkService) GetByCourse(courseID uint) ([]HomeworkResponse, error) {
	homeworks, err := s.repo.GetByCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homeworks: %w", err)
	}
	return s.toResponseDTOList(homeworks), nil
}

// GetUpcoming retrieves upcoming homework assignments
func (s *homeworkService) GetUpcoming(limit int) ([]HomeworkResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	homeworks, err := s.repo.GetUpcoming(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming homeworks: %w", err)
	}
	return s.toResponseDTOList(homeworks), nil
}

// Update updates a homework assignment
func (s *homeworkService) Update(id uint, req *UpdateHomeworkRequest) (*HomeworkResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	hw, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("homework not found: %w", err)
	}

	// Update fields
	if req.Title != "" {
		hw.Title = req.Title
	}
	if req.Description != "" {
		hw.Description = req.Description
	}
	if req.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date format: %w", err)
		}
		hw.DueDate = dueDate
	}
	if req.MaxScore != 0 {
		hw.MaxScore = req.MaxScore
	}

	// Save
	if err := s.repo.Update(hw); err != nil {
		return nil, fmt.Errorf("failed to update homework: %w", err)
	}

	return s.toResponseDTO(hw), nil
}

// Delete deletes a homework assignment
func (s *homeworkService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("homework not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete homework: %w", err)
	}

	return nil
}

// Validation methods
func (s *homeworkService) validateCreateRequest(req *CreateHomeworkRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if req.CourseID == 0 {
		return fmt.Errorf("course ID is required")
	}
	if req.DueDate == "" {
		return fmt.Errorf("due date is required")
	}
	if req.MaxScore <= 0 {
		return fmt.Errorf("max score must be greater than 0")
	}
	return nil
}

func (s *homeworkService) validateUpdateRequest(req *UpdateHomeworkRequest) error {
	if req.MaxScore != 0 && req.MaxScore <= 0 {
		return fmt.Errorf("max score must be greater than 0")
	}
	return nil
}

// DTO mapping methods
func (s *homeworkService) toResponseDTO(hw *Homework) *HomeworkResponse {
	return &HomeworkResponse{
		ID:          hw.ID,
		Title:       hw.Title,
		Description: hw.Description,
		CourseID:    hw.CourseID,
		DueDate:     hw.DueDate,
		MaxScore:    hw.MaxScore,
		CreatedAt:   hw.CreatedAt,
		UpdatedAt:   hw.UpdatedAt,
	}
}

func (s *homeworkService) toResponseDTOList(homeworks []Homework) []HomeworkResponse {
	responses := make([]HomeworkResponse, len(homeworks))
	for i, hw := range homeworks {
		responses[i] = *s.toResponseDTO(&hw)
	}
	return responses
}
