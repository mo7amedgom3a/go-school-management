package exam

import (
	"fmt"
	"strings"
	"time"
)

// ExamService defines the business logic interface
type ExamService interface {
	Create(req *CreateExamRequest) (*ExamResponse, error)
	GetByID(id uint) (*ExamResponse, error)
	GetByCourse(courseID uint) ([]ExamResponse, error)
	GetUpcoming(limit int) ([]ExamResponse, error)
	Update(id uint, req *UpdateExamRequest) (*ExamResponse, error)
	Delete(id uint) error
}

// examService implements ExamService
type examService struct {
	repo ExamRepository
}

// NewExamService creates a new exam service with DI
func NewExamService(repo ExamRepository) ExamService {
	return &examService{repo: repo}
}

// Create creates a new exam
func (s *examService) Create(req *CreateExamRequest) (*ExamResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse exam date
	examDate, err := time.Parse(time.RFC3339, req.ExamDate)
	if err != nil {
		return nil, fmt.Errorf("invalid exam date format (use RFC3339): %w", err)
	}

	// Map DTO to Model
	ex := &Exam{
		Title:    req.Title,
		CourseID: req.CourseID,
		ExamDate: examDate,
		Duration: req.Duration,
		MaxScore: req.MaxScore,
	}

	// Create via repository
	if err := s.repo.Create(ex); err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(ex), nil
}

// GetByID retrieves an exam by ID
func (s *examService) GetByID(id uint) (*ExamResponse, error) {
	ex, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}
	return s.toResponseDTO(ex), nil
}

// GetByCourse retrieves exams for a course
func (s *examService) GetByCourse(courseID uint) ([]ExamResponse, error) {
	exams, err := s.repo.GetByCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exams: %w", err)
	}
	return s.toResponseDTOList(exams), nil
}

// GetUpcoming retrieves upcoming exams
func (s *examService) GetUpcoming(limit int) ([]ExamResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	exams, err := s.repo.GetUpcoming(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming exams: %w", err)
	}
	return s.toResponseDTOList(exams), nil
}

// Update updates an exam
func (s *examService) Update(id uint, req *UpdateExamRequest) (*ExamResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	ex, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	// Update fields
	if req.Title != "" {
		ex.Title = req.Title
	}
	if req.ExamDate != "" {
		examDate, err := time.Parse(time.RFC3339, req.ExamDate)
		if err != nil {
			return nil, fmt.Errorf("invalid exam date format: %w", err)
		}
		ex.ExamDate = examDate
	}
	if req.Duration != 0 {
		ex.Duration = req.Duration
	}
	if req.MaxScore != 0 {
		ex.MaxScore = req.MaxScore
	}

	// Save
	if err := s.repo.Update(ex); err != nil {
		return nil, fmt.Errorf("failed to update exam: %w", err)
	}

	return s.toResponseDTO(ex), nil
}

// Delete deletes an exam
func (s *examService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete exam: %w", err)
	}

	return nil
}

// Validation methods
func (s *examService) validateCreateRequest(req *CreateExamRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if req.CourseID == 0 {
		return fmt.Errorf("course ID is required")
	}
	if req.ExamDate == "" {
		return fmt.Errorf("exam date is required")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if req.MaxScore <= 0 {
		return fmt.Errorf("max score must be greater than 0")
	}
	return nil
}

func (s *examService) validateUpdateRequest(req *UpdateExamRequest) error {
	if req.Duration != 0 && req.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if req.MaxScore != 0 && req.MaxScore <= 0 {
		return fmt.Errorf("max score must be greater than 0")
	}
	return nil
}

// DTO mapping methods
func (s *examService) toResponseDTO(ex *Exam) *ExamResponse {
	return &ExamResponse{
		ID:        ex.ID,
		Title:     ex.Title,
		CourseID:  ex.CourseID,
		ExamDate:  ex.ExamDate,
		Duration:  ex.Duration,
		MaxScore:  ex.MaxScore,
		CreatedAt: ex.CreatedAt,
		UpdatedAt: ex.UpdatedAt,
	}
}

func (s *examService) toResponseDTOList(exams []Exam) []ExamResponse {
	responses := make([]ExamResponse, len(exams))
	for i, ex := range exams {
		responses[i] = *s.toResponseDTO(&ex)
	}
	return responses
}
