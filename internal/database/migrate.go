package database

import (
	"log"

	"school_management/internal/modules/attendance"
	"school_management/internal/modules/course"
	"school_management/internal/modules/department"
	"school_management/internal/modules/exam"
	"school_management/internal/modules/grade"
	"school_management/internal/modules/homework"
	"school_management/internal/modules/student"
	"school_management/internal/modules/student_courses"
	"school_management/internal/modules/students_homework"
	"school_management/internal/modules/teacher"
)

// RunMigrations runs all database migrations using GORM AutoMigrate
func RunMigrations() error {
	log.Println("🔄 Running database migrations...")

	// Order matters: migrate parent tables before child tables
	err := DB.AutoMigrate(
		// Core entities (no dependencies)
		&department.Department{},

		// Entities with single dependencies
		&teacher.Teacher{}, // depends on Department
		&student.Student{}, // no dependencies
		&course.Course{},   // depends on Department and Teacher

		// Academic operations (depend on core entities)
		&attendance.Attendance{}, // depends on Student and Course
		&homework.Homework{},     // depends on Course
		&exam.Exam{},             // depends on Course
		&grade.Grade{},           // depends on Student and Exam

		// Junction tables (depend on multiple entities)
		&student_courses.StudentCourse{},     // depends on Student and Course
		&students_homework.StudentHomework{}, // depends on Student and Homework
	)

	if err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migrations completed successfully!")
	return nil
}
