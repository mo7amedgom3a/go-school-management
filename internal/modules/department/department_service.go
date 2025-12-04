package department

import (
	"fmt"
	"strings"
)

// DepartmentService defines the business logic interface
type DepartmentService interface {
	Create(req *CreateDepartmentRequest) (*DepartmentResponse, error)
	GetByID(id uint) (*DepartmentResponse, error)
	GetAll(limit, offset int) ([]DepartmentResponse, error)
	Update(id uint, req *UpdateDepartmentRequest) (*DepartmentResponse, error)
	Delete(id uint) error
	Search(name string) ([]DepartmentResponse, error)
}

// departmentService implements DepartmentService
type departmentService struct {
	repo DepartmentRepository
}

// NewDepartmentService creates a new department service with DI
func NewDepartmentService(repo DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

// Create creates a new department
func (s *departmentService) Create(req *CreateDepartmentRequest) (*DepartmentResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Map DTO to Model
	dept := &Department{
		Name:        req.Name,
		Description: req.Description,
	}

	// Create via repository
	if err := s.repo.Create(dept); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(dept), nil
}

// GetByID retrieves a department by ID
func (s *departmentService) GetByID(id uint) (*DepartmentResponse, error) {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}
	return s.toResponseDTO(dept), nil
}

// GetAll retrieves all departments with pagination
func (s *departmentService) GetAll(limit, offset int) ([]DepartmentResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	departments, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}

	return s.toResponseDTOList(departments), nil
}

// Update updates a department
func (s *departmentService) Update(id uint, req *UpdateDepartmentRequest) (*DepartmentResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	// Update fields
	if req.Name != "" {
		dept.Name = req.Name
	}
	if req.Description != "" {
		dept.Description = req.Description
	}

	// Save
	if err := s.repo.Update(dept); err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return s.toResponseDTO(dept), nil
}

// Delete deletes a department
func (s *departmentService) Delete(id uint) error {
	// Check if exists
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("department not found: %w", err)
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}

// Search searches departments by name
func (s *departmentService) Search(name string) ([]DepartmentResponse, error) {
	departments, err := s.repo.Search(name)
	if err != nil {
		return nil, fmt.Errorf("failed to search departments: %w", err)
	}
	return s.toResponseDTOList(departments), nil
}

// Validation methods
func (s *departmentService) validateCreateRequest(req *CreateDepartmentRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("department name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("department name must be less than 100 characters")
	}
	return nil
}

func (s *departmentService) validateUpdateRequest(req *UpdateDepartmentRequest) error {
	if req.Name != "" && len(req.Name) > 100 {
		return fmt.Errorf("department name must be less than 100 characters")
	}
	return nil
}

// DTO mapping methods
func (s *departmentService) toResponseDTO(dept *Department) *DepartmentResponse {
	return &DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
	}
}

func (s *departmentService) toResponseDTOList(departments []Department) []DepartmentResponse {
	responses := make([]DepartmentResponse, len(departments))
	for i, dept := range departments {
		responses[i] = *s.toResponseDTO(&dept)
	}
	return responses
}
