package course

import (
	"fmt"
	"strings"
)

// CourseService defines the business logic interface
type CourseService interface {
	Create(req *CreateCourseRequest) (*CourseResponse, error)
	GetByID(id uint) (*CourseResponse, error)
	GetByIDWithRelations(id uint) (*CourseResponse, error)
	GetAll(limit, offset int) ([]CourseResponse, error)
	GetByDepartment(deptID uint) ([]CourseResponse, error)
	Update(id uint, req *UpdateCourseRequest) (*CourseResponse, error)
	Delete(id uint) error
}

// courseService implements CourseService
type courseService struct {
	repo CourseRepository
}

// NewCourseService creates a new course service with DI
func NewCourseService(repo CourseRepository) CourseService {
	return &courseService{repo: repo}
}

// Create creates a new course
func (s *courseService) Create(req *CreateCourseRequest) (*CourseResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Map DTO to Model
	course := &Course{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Credits:      req.Credits,
		DepartmentID: req.DepartmentID,
		TeacherID:    req.TeacherID,
	}

	// Create via repository
	if err := s.repo.Create(course); err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(course), nil
}

// GetByID retrieves a course by ID
func (s *courseService) GetByID(id uint) (*CourseResponse, error) {
	course, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	return s.toResponseDTO(course), nil
}

// GetByIDWithRelations retrieves a course with relations preloaded
func (s *courseService) GetByIDWithRelations(id uint) (*CourseResponse, error) {
	course, err := s.repo.GetByIDWithRelations(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	return s.toResponseDTO(course), nil
}

// GetAll retrieves all courses with pagination
func (s *courseService) GetAll(limit, offset int) ([]CourseResponse, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	courses, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return s.toResponseDTOList(courses), nil
}

// GetByDepartment retrieves courses by department
func (s *courseService) GetByDepartment(deptID uint) ([]CourseResponse, error) {
	courses, err := s.repo.GetByDepartment(deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses by department: %w", err)
	}
	return s.toResponseDTOList(courses), nil
}

// Update updates a course
func (s *courseService) Update(id uint, req *UpdateCourseRequest) (*CourseResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	course, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Update fields
	if req.Name != "" {
		course.Name = req.Name
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Credits != 0 {
		course.Credits = req.Credits
	}
	if req.DepartmentID != 0 {
		course.DepartmentID = req.DepartmentID
	}
	if req.TeacherID != 0 {
		course.TeacherID = req.TeacherID
	}

	// Save
	if err := s.repo.Update(course); err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	return s.toResponseDTO(course), nil
}

// Delete deletes a course
func (s *courseService) Delete(id uint) error {
	// Check if exists
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// Delete
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	return nil
}

// Validation methods
func (s *courseService) validateCreateRequest(req *CreateCourseRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("course name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("course code is required")
	}
	if req.Credits < 1 || req.Credits > 6 {
		return fmt.Errorf("credits must be between 1 and 6")
	}
	if req.DepartmentID == 0 {
		return fmt.Errorf("department ID is required")
	}
	if req.TeacherID == 0 {
		return fmt.Errorf("teacher ID is required")
	}
	return nil
}

func (s *courseService) validateUpdateRequest(req *UpdateCourseRequest) error {
	if req.Credits != 0 && (req.Credits < 1 || req.Credits > 6) {
		return fmt.Errorf("credits must be between 1 and 6")
	}
	return nil
}

// DTO mapping methods
func (s *courseService) toResponseDTO(course *Course) *CourseResponse {
	return &CourseResponse{
		ID:           course.ID,
		Name:         course.Name,
		Code:         course.Code,
		Description:  course.Description,
		Credits:      course.Credits,
		DepartmentID: course.DepartmentID,
		TeacherID:    course.TeacherID,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}
}

func (s *courseService) toResponseDTOList(courses []Course) []CourseResponse {
	responses := make([]CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = *s.toResponseDTO(&course)
	}
	return responses
}
