package student

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// StudentService defines the business logic interface
type StudentService interface {
	Create(req *CreateStudentRequest) (*StudentResponse, error)
	GetByID(id uint) (*StudentResponse, error)
	GetAll(limit, offset int) ([]StudentResponse, error)
	Update(id uint, req *UpdateStudentRequest) (*StudentResponse, error)
	Delete(id uint) error
	Search(query string, limit int) ([]StudentResponse, error)
}

// studentService implements StudentService
type studentService struct {
	repo StudentRepository
}

// NewStudentService creates a new student service with DI
func NewStudentService(repo StudentRepository) StudentService {
	return &studentService{repo: repo}
}

// Create creates a new student
func (s *studentService) Create(req *CreateStudentRequest) (*StudentResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse dates
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth format (use YYYY-MM-DD): %w", err)
	}

	enrollDate, err := time.Parse("2006-01-02", req.EnrollmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid enrollment date format (use YYYY-MM-DD): %w", err)
	}

	// Map DTO to Model
	student := &Student{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DateOfBirth:    dob,
		EnrollmentDate: enrollDate,
	}

	// Create via repository
	if err := s.repo.Create(student); err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(student), nil
}

// GetByID retrieves a student by ID
func (s *studentService) GetByID(id uint) (*StudentResponse, error) {
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("student not found: %w", err)
	}
	return s.toResponseDTO(student), nil
}

// GetAll retrieves all students with pagination
func (s *studentService) GetAll(limit, offset int) ([]StudentResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	students, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}

	return s.toResponseDTOList(students), nil
}

// Update updates a student
func (s *studentService) Update(id uint, req *UpdateStudentRequest) (*StudentResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("student not found: %w", err)
	}

	// Update fields
	if req.FirstName != "" {
		student.FirstName = req.FirstName
	}
	if req.LastName != "" {
		student.LastName = req.LastName
	}
	if req.Email != "" {
		student.Email = req.Email
	}
	if req.Phone != "" {
		student.Phone = req.Phone
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
		student.DateOfBirth = dob
	}

	// Save
	if err := s.repo.Update(student); err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	return s.toResponseDTO(student), nil
}

// Delete deletes a student
func (s *studentService) Delete(id uint) error {
	// Check if exists
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	return nil
}

// Search searches students
func (s *studentService) Search(query string, limit int) ([]StudentResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	students, err := s.repo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search students: %w", err)
	}

	return s.toResponseDTOList(students), nil
}

// Validation methods
func (s *studentService) validateCreateRequest(req *CreateStudentRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !s.isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.DateOfBirth == "" {
		return fmt.Errorf("date of birth is required")
	}
	if req.EnrollmentDate == "" {
		return fmt.Errorf("enrollment date is required")
	}
	return nil
}

func (s *studentService) validateUpdateRequest(req *UpdateStudentRequest) error {
	if req.Email != "" && !s.isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (s *studentService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// DTO mapping methods
func (s *studentService) toResponseDTO(student *Student) *StudentResponse {
	return &StudentResponse{
		ID:             student.ID,
		FirstName:      student.FirstName,
		LastName:       student.LastName,
		Email:          student.Email,
		Phone:          student.Phone,
		DateOfBirth:    student.DateOfBirth,
		EnrollmentDate: student.EnrollmentDate,
		CreatedAt:      student.CreatedAt,
		UpdatedAt:      student.UpdatedAt,
	}
}

func (s *studentService) toResponseDTOList(students []Student) []StudentResponse {
	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = *s.toResponseDTO(&student)
	}
	return responses
}
