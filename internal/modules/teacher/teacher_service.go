package teacher

import (
	"fmt"
	"regexp"
	"strings"
)

// TeacherService defines the business logic interface
type TeacherService interface {
	Create(req *CreateTeacherRequest) (*TeacherResponse, error)
	GetByID(id uint) (*TeacherResponse, error)
	GetByIDWithDepartment(id uint) (*TeacherResponse, error)
	GetAll(limit, offset int) ([]TeacherResponse, error)
	GetByDepartment(deptID uint) ([]TeacherResponse, error)
	Update(id uint, req *UpdateTeacherRequest) (*TeacherResponse, error)
	Delete(id uint) error
}

// teacherService implements TeacherService
type teacherService struct {
	repo TeacherRepository
}

// NewTeacherService creates a new teacher service with DI
func NewTeacherService(repo TeacherRepository) TeacherService {
	return &teacherService{repo: repo}
}

// Create creates a new teacher
func (s *teacherService) Create(req *CreateTeacherRequest) (*TeacherResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Map DTO to Model
	teacher := &Teacher{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
	}

	// Create via repository
	if err := s.repo.Create(teacher); err != nil {
		return nil, fmt.Errorf("failed to create teacher: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(teacher), nil
}

// GetByID retrieves a teacher by ID
func (s *teacherService) GetByID(id uint) (*TeacherResponse, error) {
	teacher, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("teacher not found: %w", err)
	}
	return s.toResponseDTO(teacher), nil
}

// GetByIDWithDepartment retrieves a teacher with department preloaded
func (s *teacherService) GetByIDWithDepartment(id uint) (*TeacherResponse, error) {
	teacher, err := s.repo.GetByIDWithDepartment(id)
	if err != nil {
		return nil, fmt.Errorf("teacher not found: %w", err)
	}
	return s.toResponseDTO(teacher), nil
}

// GetAll retrieves all teachers with pagination
func (s *teacherService) GetAll(limit, offset int) ([]TeacherResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	teachers, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers: %w", err)
	}

	return s.toResponseDTOList(teachers), nil
}

// GetByDepartment retrieves teachers by department
func (s *teacherService) GetByDepartment(deptID uint) ([]TeacherResponse, error) {
	teachers, err := s.repo.GetByDepartment(deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers by department: %w", err)
	}
	return s.toResponseDTOList(teachers), nil
}

// Update updates a teacher
func (s *teacherService) Update(id uint, req *UpdateTeacherRequest) (*TeacherResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	teacher, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("teacher not found: %w", err)
	}

	// Update fields
	if req.FirstName != "" {
		teacher.FirstName = req.FirstName
	}
	if req.LastName != "" {
		teacher.LastName = req.LastName
	}
	if req.Email != "" {
		teacher.Email = req.Email
	}
	if req.Phone != "" {
		teacher.Phone = req.Phone
	}
	if req.DepartmentID != 0 {
		teacher.DepartmentID = req.DepartmentID
	}

	// Save
	if err := s.repo.Update(teacher); err != nil {
		return nil, fmt.Errorf("failed to update teacher: %w", err)
	}

	return s.toResponseDTO(teacher), nil
}

// Delete deletes a teacher
func (s *teacherService) Delete(id uint) error {
	// Check if exists
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("teacher not found: %w", err)
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}

	return nil
}

// Validation methods
func (s *teacherService) validateCreateRequest(req *CreateTeacherRequest) error {
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
	if req.DepartmentID == 0 {
		return fmt.Errorf("department ID is required")
	}
	return nil
}

func (s *teacherService) validateUpdateRequest(req *UpdateTeacherRequest) error {
	if req.Email != "" && !s.isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (s *teacherService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// DTO mapping methods
func (s *teacherService) toResponseDTO(teacher *Teacher) *TeacherResponse {
	return &TeacherResponse{
		ID:           teacher.ID,
		FirstName:    teacher.FirstName,
		LastName:     teacher.LastName,
		Email:        teacher.Email,
		Phone:        teacher.Phone,
		DepartmentID: teacher.DepartmentID,
		CreatedAt:    teacher.CreatedAt,
		UpdatedAt:    teacher.UpdatedAt,
	}
}

func (s *teacherService) toResponseDTOList(teachers []Teacher) []TeacherResponse {
	responses := make([]TeacherResponse, len(teachers))
	for i, teacher := range teachers {
		responses[i] = *s.toResponseDTO(&teacher)
	}
	return responses
}
