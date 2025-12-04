package attendance

import (
	"fmt"
	"time"
)

// AttendanceService defines the business logic interface
type AttendanceService interface {
	Create(req *CreateAttendanceRequest) (*AttendanceResponse, error)
	GetByID(id uint) (*AttendanceResponse, error)
	GetByStudent(studentID uint) ([]AttendanceResponse, error)
	GetByCourse(courseID uint) ([]AttendanceResponse, error)
	Update(id uint, req *UpdateAttendanceRequest) (*AttendanceResponse, error)
	Delete(id uint) error
}

// attendanceService implements AttendanceService
type attendanceService struct {
	repo AttendanceRepository
}

// NewAttendanceService creates a new attendance service with DI
func NewAttendanceService(repo AttendanceRepository) AttendanceService {
	return &attendanceService{repo: repo}
}

// Create creates a new attendance record
func (s *attendanceService) Create(req *CreateAttendanceRequest) (*AttendanceResponse, error) {
	// Validate
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
	}

	// Map DTO to Model
	att := &Attendance{
		StudentID: req.StudentID,
		CourseID:  req.CourseID,
		Date:      date,
		Status:    AttendanceStatus(req.Status),
	}

	// Create via repository
	if err := s.repo.Create(att); err != nil {
		return nil, fmt.Errorf("failed to create attendance: %w", err)
	}

	// Map Model to Response DTO
	return s.toResponseDTO(att), nil
}

// GetByID retrieves an attendance record by ID
func (s *attendanceService) GetByID(id uint) (*AttendanceResponse, error) {
	att, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}
	return s.toResponseDTO(att), nil
}

// GetByStudent retrieves attendance records for a student
func (s *attendanceService) GetByStudent(studentID uint) ([]AttendanceResponse, error) {
	attendances, err := s.repo.GetByStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}
	return s.toResponseDTOList(attendances), nil
}

// GetByCourse retrieves attendance records for a course
func (s *attendanceService) GetByCourse(courseID uint) ([]AttendanceResponse, error) {
	attendances, err := s.repo.GetByCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}
	return s.toResponseDTOList(attendances), nil
}

// Update updates an attendance record
func (s *attendanceService) Update(id uint, req *UpdateAttendanceRequest) (*AttendanceResponse, error) {
	// Validate
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Get existing
	att, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}

	// Update fields
	if req.Status != "" {
		att.Status = AttendanceStatus(req.Status)
	}

	// Save
	if err := s.repo.Update(att); err != nil {
		return nil, fmt.Errorf("failed to update attendance: %w", err)
	}

	return s.toResponseDTO(att), nil
}

// Delete deletes an attendance record
func (s *attendanceService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("attendance not found: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete attendance: %w", err)
	}

	return nil
}

// Validation methods
func (s *attendanceService) validateCreateRequest(req *CreateAttendanceRequest) error {
	if req.StudentID == 0 {
		return fmt.Errorf("student ID is required")
	}
	if req.CourseID == 0 {
		return fmt.Errorf("course ID is required")
	}
	if req.Date == "" {
		return fmt.Errorf("date is required")
	}
	if !s.isValidStatus(req.Status) {
		return fmt.Errorf("invalid status (must be: present, absent, or late)")
	}
	return nil
}

func (s *attendanceService) validateUpdateRequest(req *UpdateAttendanceRequest) error {
	if req.Status != "" && !s.isValidStatus(req.Status) {
		return fmt.Errorf("invalid status (must be: present, absent, or late)")
	}
	return nil
}

func (s *attendanceService) isValidStatus(status string) bool {
	return status == "present" || status == "absent" || status == "late"
}

// DTO mapping methods
func (s *attendanceService) toResponseDTO(att *Attendance) *AttendanceResponse {
	return &AttendanceResponse{
		ID:        att.ID,
		StudentID: att.StudentID,
		CourseID:  att.CourseID,
		Date:      att.Date,
		Status:    string(att.Status),
		CreatedAt: att.CreatedAt,
		UpdatedAt: att.UpdatedAt,
	}
}

func (s *attendanceService) toResponseDTOList(attendances []Attendance) []AttendanceResponse {
	responses := make([]AttendanceResponse, len(attendances))
	for i, att := range attendances {
		responses[i] = *s.toResponseDTO(&att)
	}
	return responses
}
