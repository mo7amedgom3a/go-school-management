package students_homework

import (
	"fmt"
	"time"
)

// StudentHomeworkService defines the business logic interface
type StudentHomeworkService interface {
	Submit(req *SubmitHomeworkRequest) (*StudentHomeworkResponse, error)
	Grade(req *GradeHomeworkRequest) (*StudentHomeworkResponse, error)
	GetByID(id uint) (*StudentHomeworkResponse, error)
	GetByStudent(studentID uint) ([]StudentHomeworkResponse, error)
	GetByHomework(homeworkID uint) ([]StudentHomeworkResponse, error)
	GetPendingByStudent(studentID uint) ([]StudentHomeworkResponse, error)
	Delete(id uint) error
}

// studentHomeworkService implements StudentHomeworkService
type studentHomeworkService struct {
	repo StudentHomeworkRepository
}

// NewStudentHomeworkService creates a new student homework service with DI
func NewStudentHomeworkService(repo StudentHomeworkRepository) StudentHomeworkService {
	return &studentHomeworkService{repo: repo}
}

// Submit submits homework
func (s *studentHomeworkService) Submit(req *SubmitHomeworkRequest) (*StudentHomeworkResponse, error) {
	// Validate
	if err := s.validateSubmitRequest(req); err != nil {
		return nil, err
	}

	// Check if already submitted
	existing, _ := s.repo.GetByStudentAndHomework(req.StudentID, req.HomeworkID)
	if existing != nil {
		return nil, fmt.Errorf("homework already submitted")
	}

	// Map DTO to Model
	now := time.Now()
	submission := &StudentHomework{
		StudentID:      req.StudentID,
		HomeworkID:     req.HomeworkID,
		SubmissionDate: &now,
		Status:         HomeworkSubmitted,
	}

	// Create via repository
	if err := s.repo.Create(submission); err != nil {
		return nil, fmt.Errorf("failed to submit homework: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(submission), nil
}

// Grade grades a homework submission
func (s *studentHomeworkService) Grade(req *GradeHomeworkRequest) (*StudentHomeworkResponse, error) {
	// Validate
	if err := s.validateGradeRequest(req); err != nil {
		return nil, err
	}

	// Get existing submission
	submission, err := s.repo.GetByStudentAndHomework(req.StudentID, req.HomeworkID)
	if err != nil {
		return nil, fmt.Errorf("submission not found: %w", err)
	}

	// Update grade
	submission.Score = &req.Score
	submission.Status = HomeworkGraded

	// Save
	if err := s.repo.Update(submission); err != nil {
		return nil, fmt.Errorf("failed to grade homework: %w", err)
	}

	return s.toResponseDTO(submission), nil
}

// GetByID retrieves a submission by ID
func (s *studentHomeworkService) GetByID(id uint) (*StudentHomeworkResponse, error) {
	submission, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("submission not found: %w", err)
	}
	return s.toResponseDTO(submission), nil
}

// GetByStudent retrieves submissions for a student
func (s *studentHomeworkService) GetByStudent(studentID uint) ([]StudentHomeworkResponse, error) {
	submissions, err := s.repo.GetByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions: %w", err)
	}
	return s.toResponseDTOList(submissions), nil
}

// GetByHomework retrieves submissions for a homework
func (s *studentHomeworkService) GetByHomework(homeworkID uint) ([]StudentHomeworkResponse, error) {
	submissions, err := s.repo.GetByHomework(homeworkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get submissions: %w", err)
	}
	return s.toResponseDTOList(submissions), nil
}

// GetPendingByStudent retrieves pending submissions for a student
func (s *studentHomeworkService) GetPendingByStudent(studentID uint) ([]StudentHomeworkResponse, error) {
	submissions, err := s.repo.GetPendingByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending submissions: %w", err)
	}
	return s.toResponseDTOList(submissions), nil
}

// Delete deletes a submission
func (s *studentHomeworkService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete submission: %w", err)
	}

	return nil
}

// Validation methods
func (s *studentHomeworkService) validateSubmitRequest(req *SubmitHomeworkRequest) error {
	if req.StudentID == 0 {
		return fmt.Errorf("student ID is required")
	}
	if req.HomeworkID == 0 {
		return fmt.Errorf("homework ID is required")
	}
	return nil
}

func (s *studentHomeworkService) validateGradeRequest(req *GradeHomeworkRequest) error {
	if req.StudentID == 0 {
		return fmt.Errorf("student ID is required")
	}
	if req.HomeworkID == 0 {
		return fmt.Errorf("homework ID is required")
	}
	if req.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}
	return nil
}

// DTO mapping methods
func (s *studentHomeworkService) toResponseDTO(submission *StudentHomework) *StudentHomeworkResponse {
	resp := &StudentHomeworkResponse{
		ID:         submission.ID,
		StudentID:  submission.StudentID,
		HomeworkID: submission.HomeworkID,
		Status:     string(submission.Status),
		CreatedAt:  submission.CreatedAt,
		UpdatedAt:  submission.UpdatedAt,
	}

	if submission.SubmissionDate != nil {
		submissionDate := submission.SubmissionDate.Format(time.RFC3339)
		resp.SubmissionDate = &submissionDate
	}

	if submission.Score != nil {
		resp.Score = submission.Score
	}

	return resp
}

func (s *studentHomeworkService) toResponseDTOList(submissions []StudentHomework) []StudentHomeworkResponse {
	responses := make([]StudentHomeworkResponse, len(submissions))
	for i, submission := range submissions {
		responses[i] = *s.toResponseDTO(&submission)
	}
	return responses
}
