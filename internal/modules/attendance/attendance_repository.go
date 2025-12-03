package attendance

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AttendanceRepository defines the interface for attendance data access
type AttendanceRepository interface {
	Create(attendance *Attendance) error
	GetByID(id uint) (*Attendance, error)
	GetByIDWithRelations(id uint) (*Attendance, error)
	GetAll(limit, offset int) ([]Attendance, error)
	GetByStudent(studentID uint) ([]Attendance, error)
	GetByCourse(courseID uint) ([]Attendance, error)
	GetByDateRange(start, end time.Time) ([]Attendance, error)
	GetByStudentAndCourse(studentID, courseID uint) ([]Attendance, error)
	Update(attendance *Attendance) error
	Delete(id uint) error
}

// attendanceRepository implements AttendanceRepository
type attendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new attendance repository with dependency injection
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

// Create creates a new attendance record
func (r *attendanceRepository) Create(attendance *Attendance) error {
	if err := r.db.Create(attendance).Error; err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}
	return nil
}

// GetByID retrieves an attendance record by ID
func (r *attendanceRepository) GetByID(id uint) (*Attendance, error) {
	var attendance Attendance
	if err := r.db.First(&attendance, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}
	return &attendance, nil
}

// GetByIDWithRelations retrieves an attendance record with student and course preloaded
func (r *attendanceRepository) GetByIDWithRelations(id uint) (*Attendance, error) {
	var attendance Attendance
	if err := r.db.Preload("Student").Preload("Course").First(&attendance, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendance with relations: %w", err)
	}
	return &attendance, nil
}

// GetAll retrieves all attendance records with pagination
func (r *attendanceRepository) GetAll(limit, offset int) ([]Attendance, error) {
	var attendances []Attendance
	if err := r.db.Limit(limit).Offset(offset).Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}
	return attendances, nil
}

// GetByStudent retrieves all attendance records for a student
func (r *attendanceRepository) GetByStudent(studentID uint) ([]Attendance, error) {
	var attendances []Attendance
	if err := r.db.Where("student_id = ?", studentID).Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendances by student: %w", err)
	}
	return attendances, nil
}

// GetByCourse retrieves all attendance records for a course
func (r *attendanceRepository) GetByCourse(courseID uint) ([]Attendance, error) {
	var attendances []Attendance
	if err := r.db.Where("course_id = ?", courseID).Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendances by course: %w", err)
	}
	return attendances, nil
}

// GetByDateRange retrieves attendance records within a date range
func (r *attendanceRepository) GetByDateRange(start, end time.Time) ([]Attendance, error) {
	var attendances []Attendance
	if err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendances by date range: %w", err)
	}
	return attendances, nil
}

// GetByStudentAndCourse retrieves attendance records for a specific student in a specific course
func (r *attendanceRepository) GetByStudentAndCourse(studentID, courseID uint) ([]Attendance, error) {
	var attendances []Attendance
	if err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).
		Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to get attendances by student and course: %w", err)
	}
	return attendances, nil
}

// Update updates an attendance record
func (r *attendanceRepository) Update(attendance *Attendance) error {
	if err := r.db.Save(attendance).Error; err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}
	return nil
}

// Delete soft deletes an attendance record
func (r *attendanceRepository) Delete(id uint) error {
	if err := r.db.Delete(&Attendance{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete attendance: %w", err)
	}
	return nil
}
