package grade

import (
	"fmt"
)

// GradeService defines the business logic interface
type GradeService interface {
	Create(req *CreateGradeRequest) (*GradeResponse, error)
	GetByID(id uint) (*GradeResponse, error)
	GetByStudent(studentID uint) ([]GradeResponse, error)
	GetByExam(examID uint) ([]GradeResponse, error)
	GetStudentAverage(studentID uint) (float64, error)
	Update(id uint, req *UpdateGradeRequest) (*GradeResponse, error)
	Delete(id uint) error
}

// gradeService implements GradeService
type gradeService struct {
	repo GradeRepository
}

// NewGradeService creates a new grade service with DI
func NewGradeService(repo GradeRepository) GradeService {
	return &gradeService{repo: repo}
}

// Create creates a new grade
func (s *gradeService) Create(req *CreateGradeRequest) (*GradeResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Map DTO to Model
	gr := &Grade{
		StudentID: req.StudentID,
		ExamID:    req.ExamID,
		Score:     req.Score,
	}

	// Create via repository
	if err := s.repo.Create(gr); err != nil {
		return nil, fmt.Errorf("failed to create grade: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(gr), nil
}

// GetByID retrieves a grade by ID
func (s *gradeService) GetByID(id uint) (*GradeResponse, error) {
	gr, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("grade not found: %w", err)
	}
	return s.toResponseDTO(gr), nil
}

// GetByStudent retrieves grades for a student
func (s *gradeService) GetByStudent(studentID uint) ([]GradeResponse, error) {
	grades, err := s.repo.GetByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get grades: %w", err)
	}
	return s.toResponseDTOList(grades), nil
}

// GetByExam retrieves grades for an exam
func (s *gradeService) GetByExam(examID uint) ([]GradeResponse, error) {
	grades, err := s.repo.GetByExam(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get grades: %w", err)
	}
	return s.toResponseDTOList(grades), nil
}

// GetStudentAverage calculates a student's average grade
func (s *gradeService) GetStudentAverage(studentID uint) (float64, error) {
	avg, err := s.repo.GetStudentAverage(studentID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate average: %w", err)
	}
	return avg, nil
}

// Update updates a grade
func (s *gradeService) Update(id uint, req *UpdateGradeRequest) (*GradeResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	gr, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("grade not found: %w", err)
	}

	// Update fields
	if req.Score != 0 {
		gr.Score = req.Score
	}

	// Save
	if err := s.repo.Update(gr); err != nil {
		return nil, fmt.Errorf("failed to update grade: %w", err)
	}

	return s.toResponseDTO(gr), nil
}

// Delete deletes a grade
func (s *gradeService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("grade not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete grade: %w", err)
	}

	return nil
}

// Validation methods
func (s *gradeService) validateCreateRequest(req *CreateGradeRequest) error {
	if req.StudentID == 0 {
		return fmt.Errorf("student ID is required")
	}
	if req.ExamID == 0 {
		return fmt.Errorf("exam ID is required")
	}
	if req.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}
	return nil
}

func (s *gradeService) validateUpdateRequest(req *UpdateGradeRequest) error {
	if req.Score != 0 && req.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}
	return nil
}

// DTO mapping methods
func (s *gradeService) toResponseDTO(gr *Grade) *GradeResponse {
	return &GradeResponse{
		ID:        gr.ID,
		StudentID: gr.StudentID,
		ExamID:    gr.ExamID,
		Score:     gr.Score,
		CreatedAt: gr.CreatedAt,
		UpdatedAt: gr.UpdatedAt,
	}
}

func (s *gradeService) toResponseDTOList(grades []Grade) []GradeResponse {
	responses := make([]GradeResponse, len(grades))
	for i, gr := range grades {
		responses[i] = *s.toResponseDTO(&gr)
	}
	return responses
}
