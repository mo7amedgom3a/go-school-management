package student_courses

import (
	"fmt"
	"time"
)

// StudentCourseService defines the business logic interface
type StudentCourseService interface {
	Enroll(req *EnrollStudentRequest) (*StudentCourseResponse, error)
	GetByID(id uint) (*StudentCourseResponse, error)
	GetByStudent(studentID uint) ([]StudentCourseResponse, error)
	GetByCourse(courseID uint) ([]StudentCourseResponse, error)
	Unenroll(id uint) error
}

// studentCourseService implements StudentCourseService
type studentCourseService struct {
	repo StudentCourseRepository
}

// NewStudentCourseService creates a new student course service with DI
func NewStudentCourseService(repo StudentCourseRepository) StudentCourseService {
	return &studentCourseService{repo: repo}
}

// Enroll enrolls a student in a course
func (s *studentCourseService) Enroll(req *EnrollStudentRequest) (*StudentCourseResponse, error) {
	// Validate
	if err := s.validateEnrollRequest(req); err != nil {
		return nil, err
	}

	// Check if already enrolled
	existing, _ := s.repo.GetByStudentAndCourse(req.StudentID, req.CourseID)
	if existing != nil {
		return nil, fmt.Errorf("student already enrolled in this course")
	}

	// Parse enrollment date
	enrollDate, err := time.Parse("2006-01-02", req.EnrollmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid enrollment date format (use YYYY-MM-DD): %w", err)
	}

	// Map DTO to Model
	enrollment := &StudentCourse{
		StudentID:      req.StudentID,
		CourseID:       req.CourseID,
		EnrollmentDate: enrollDate,
	}

	// Create via repository
	if err := s.repo.Create(enrollment); err != nil {
		return nil, fmt.Errorf("failed to enroll student: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(enrollment), nil
}

// GetByID retrieves an enrollment by ID
func (s *studentCourseService) GetByID(id uint) (*StudentCourseResponse, error) {
	enrollment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("enrollment not found: %w", err)
	}
	return s.toResponseDTO(enrollment), nil
}

// GetByStudent retrieves enrollments for a student
func (s *studentCourseService) GetByStudent(studentID uint) ([]StudentCourseResponse, error) {
	enrollments, err := s.repo.GetByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	return s.toResponseDTOList(enrollments), nil
}

// GetByCourse retrieves enrollments for a course
func (s *studentCourseService) GetByCourse(courseID uint) ([]StudentCourseResponse, error) {
	enrollments, err := s.repo.GetByCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	return s.toResponseDTOList(enrollments), nil
}

// Unenroll removes a student from a course
func (s *studentCourseService) Unenroll(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to unenroll student: %w", err)
	}

	return nil
}

// Validation methods
func (s *studentCourseService) validateEnrollRequest(req *EnrollStudentRequest) error {
	if req.StudentID == 0 {
		return fmt.Errorf("student ID is required")
	}
	if req.CourseID == 0 {
		return fmt.Errorf("course ID is required")
	}
	if req.EnrollmentDate == "" {
		return fmt.Errorf("enrollment date is required")
	}
	return nil
}

// DTO mapping methods
func (s *studentCourseService) toResponseDTO(enrollment *StudentCourse) *StudentCourseResponse {
	return &StudentCourseResponse{
		ID:             enrollment.ID,
		StudentID:      enrollment.StudentID,
		CourseID:       enrollment.CourseID,
		EnrollmentDate: enrollment.EnrollmentDate,
		CreatedAt:      enrollment.CreatedAt,
		UpdatedAt:      enrollment.UpdatedAt,
	}
}

func (s *studentCourseService) toResponseDTOList(enrollments []StudentCourse) []StudentCourseResponse {
	responses := make([]StudentCourseResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toResponseDTO(&enrollment)
	}
	return responses
}
