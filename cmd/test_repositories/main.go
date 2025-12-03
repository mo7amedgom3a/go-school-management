package main

import (
	"fmt"
	"log"
	"time"

	"school_management/internal/config"
	"school_management/internal/database"
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

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database.CreateDatabaseIfNotExists(cfg)
	database.ConnectDB(cfg)
	database.RunMigrations()

	fmt.Println("\n🧪 Testing Repository Layer with Dependency Injection...\n")

	// Initialize repositories with DI
	deptRepo := department.NewDepartmentRepository(database.DB)
	teacherRepo := teacher.NewTeacherRepository(database.DB)
	studentRepo := student.NewStudentRepository(database.DB)
	courseRepo := course.NewCourseRepository(database.DB)
	attendanceRepo := attendance.NewAttendanceRepository(database.DB)
	homeworkRepo := homework.NewHomeworkRepository(database.DB)
	examRepo := exam.NewExamRepository(database.DB)
	gradeRepo := grade.NewGradeRepository(database.DB)
	enrollmentRepo := student_courses.NewStudentCourseRepository(database.DB)
	submissionRepo := students_homework.NewStudentHomeworkRepository(database.DB)

	// Test Department Repository
	testDepartmentRepository(deptRepo)

	// Test Teacher Repository
	testTeacherRepository(teacherRepo, deptRepo)

	// Test Student Repository
	testStudentRepository(studentRepo)

	// Test Course Repository
	testCourseRepository(courseRepo, deptRepo, teacherRepo)

	// Test Enrollment Repository
	testEnrollmentRepository(enrollmentRepo, studentRepo, courseRepo)

	// Test Attendance Repository
	testAttendanceRepository(attendanceRepo, studentRepo, courseRepo)

	// Test Homework Repository
	testHomeworkRepository(homeworkRepo, courseRepo)

	// Test Exam Repository
	testExamRepository(examRepo, courseRepo)

	// Test Grade Repository
	testGradeRepository(gradeRepo, studentRepo, examRepo)

	// Test Submission Repository
	testSubmissionRepository(submissionRepo, studentRepo, homeworkRepo)

	fmt.Println("\n✅ All repository tests completed successfully!")
}

func testDepartmentRepository(repo department.DepartmentRepository) {
	fmt.Println("📚 Testing Department Repository...")

	// Create
	dept := &department.Department{
		Name:        "Mathematics",
		Description: "Department of Mathematics",
	}
	if err := repo.Create(dept); err != nil {
		log.Printf("Error creating department: %v", err)
		return
	}
	fmt.Printf("✓ Created department (ID: %d)\n", dept.ID)

	// Get by ID
	retrieved, err := repo.GetByID(dept.ID)
	if err != nil {
		log.Printf("Error getting department: %v", err)
		return
	}
	fmt.Printf("✓ Retrieved department: %s\n", retrieved.Name)

	// Update
	retrieved.Description = "Updated Mathematics Department"
	if err := repo.Update(retrieved); err != nil {
		log.Printf("Error updating department: %v", err)
		return
	}
	fmt.Println("✓ Updated department")

	// Get all with pagination
	departments, err := repo.GetAll(10, 0)
	if err != nil {
		log.Printf("Error getting all departments: %v", err)
		return
	}
	fmt.Printf("✓ Retrieved %d departments\n", len(departments))

	// Search
	results, err := repo.Search("Math")
	if err != nil {
		log.Printf("Error searching departments: %v", err)
		return
	}
	fmt.Printf("✓ Search found %d departments\n", len(results))

	fmt.Println()
}

func testTeacherRepository(repo teacher.TeacherRepository, deptRepo department.DepartmentRepository) {
	fmt.Println("👨‍🏫 Testing Teacher Repository...")

	// Get first department
	departments, _ := deptRepo.GetAll(1, 0)
	if len(departments) == 0 {
		fmt.Println("⚠️  No departments found, skipping teacher tests")
		return
	}

	// Create
	t := &teacher.Teacher{
		FirstName:    "Jane",
		LastName:     "Doe",
		Email:        "jane.doe@school.com",
		Phone:        "+1234567890",
		DepartmentID: departments[0].ID,
	}
	if err := repo.Create(t); err != nil {
		log.Printf("Error creating teacher: %v", err)
		return
	}
	fmt.Printf("✓ Created teacher (ID: %d)\n", t.ID)

	// Get with department (Preload)
	withDept, err := repo.GetByIDWithDepartment(t.ID)
	if err != nil {
		log.Printf("Error getting teacher with department: %v", err)
		return
	}
	fmt.Printf("✓ Retrieved teacher with department preloaded\n")
	_ = withDept

	// Get by department
	teachers, err := repo.GetByDepartment(departments[0].ID)
	if err != nil {
		log.Printf("Error getting teachers by department: %v", err)
		return
	}
	fmt.Printf("✓ Found %d teachers in department\n", len(teachers))

	// Get by email
	byEmail, err := repo.GetByEmail(t.Email)
	if err != nil {
		log.Printf("Error getting teacher by email: %v", err)
		return
	}
	fmt.Printf("✓ Found teacher by email: %s %s\n", byEmail.FirstName, byEmail.LastName)

	fmt.Println()
}

func testStudentRepository(repo student.StudentRepository) {
	fmt.Println("👨‍🎓 Testing Student Repository...")

	// Create
	dob, _ := time.Parse("2006-01-02", "2002-03-15")
	s := &student.Student{
		FirstName:      "John",
		LastName:       "Student",
		Email:          "john.student@school.com",
		Phone:          "+9876543210",
		DateOfBirth:    dob,
		EnrollmentDate: time.Now(),
	}
	if err := repo.Create(s); err != nil {
		log.Printf("Error creating student: %v", err)
		return
	}
	fmt.Printf("✓ Created student (ID: %d)\n", s.ID)

	// Search
	results, err := repo.Search("John", 10)
	if err != nil {
		log.Printf("Error searching students: %v", err)
		return
	}
	fmt.Printf("✓ Search found %d students\n", len(results))

	// Get enrolled before
	enrolled, err := repo.GetEnrolledBefore(time.Now().AddDate(0, 1, 0))
	if err != nil {
		log.Printf("Error getting enrolled students: %v", err)
		return
	}
	fmt.Printf("✓ Found %d students enrolled before date\n", len(enrolled))

	fmt.Println()
}

func testCourseRepository(repo course.CourseRepository, deptRepo department.DepartmentRepository, teacherRepo teacher.TeacherRepository) {
	fmt.Println("📖 Testing Course Repository...")

	// Get dependencies
	departments, _ := deptRepo.GetAll(1, 0)
	teachers, _ := teacherRepo.GetAll(1, 0)
	if len(departments) == 0 || len(teachers) == 0 {
		fmt.Println("⚠️  Missing dependencies, skipping course tests")
		return
	}

	// Create
	c := &course.Course{
		Name:         "Calculus I",
		Code:         "MATH101",
		Description:  "Introduction to Calculus",
		Credits:      4,
		DepartmentID: departments[0].ID,
		TeacherID:    teachers[0].ID,
	}
	if err := repo.Create(c); err != nil {
		log.Printf("Error creating course: %v", err)
		return
	}
	fmt.Printf("✓ Created course (ID: %d)\n", c.ID)

	// Get with relations (Preload)
	withRelations, err := repo.GetByIDWithRelations(c.ID)
	if err != nil {
		log.Printf("Error getting course with relations: %v", err)
		return
	}
	fmt.Printf("✓ Retrieved course with department and teacher preloaded\n")
	_ = withRelations

	// Get by code
	byCode, err := repo.GetByCode(c.Code)
	if err != nil {
		log.Printf("Error getting course by code: %v", err)
		return
	}
	fmt.Printf("✓ Found course by code: %s\n", byCode.Name)

	fmt.Println()
}

func testEnrollmentRepository(repo student_courses.StudentCourseRepository, studentRepo student.StudentRepository, courseRepo course.CourseRepository) {
	fmt.Println("📝 Testing Enrollment Repository...")

	students, _ := studentRepo.GetAll(1, 0)
	courses, _ := courseRepo.GetAll(1, 0)
	if len(students) == 0 || len(courses) == 0 {
		fmt.Println("⚠️  Missing dependencies, skipping enrollment tests")
		return
	}

	// Create enrollment
	enrollment := &student_courses.StudentCourse{
		StudentID:      students[0].ID,
		CourseID:       courses[0].ID,
		EnrollmentDate: time.Now(),
	}
	if err := repo.Create(enrollment); err != nil {
		log.Printf("Error creating enrollment: %v", err)
		return
	}
	fmt.Printf("✓ Created enrollment (ID: %d)\n", enrollment.ID)

	// Get by student
	studentEnrollments, err := repo.GetByStudent(students[0].ID)
	if err != nil {
		log.Printf("Error getting enrollments by student: %v", err)
		return
	}
	fmt.Printf("✓ Student has %d enrollments\n", len(studentEnrollments))

	fmt.Println()
}

func testAttendanceRepository(repo attendance.AttendanceRepository, studentRepo student.StudentRepository, courseRepo course.CourseRepository) {
	fmt.Println("✅ Testing Attendance Repository...")

	students, _ := studentRepo.GetAll(1, 0)
	courses, _ := courseRepo.GetAll(1, 0)
	if len(students) == 0 || len(courses) == 0 {
		fmt.Println("⚠️  Missing dependencies, skipping attendance tests")
		return
	}

	// Create attendance
	att := &attendance.Attendance{
		StudentID: students[0].ID,
		CourseID:  courses[0].ID,
		Date:      time.Now(),
		Status:    attendance.AttendancePresent,
	}
	if err := repo.Create(att); err != nil {
		log.Printf("Error creating attendance: %v", err)
		return
	}
	fmt.Printf("✓ Created attendance record (ID: %d)\n", att.ID)

	// Get by date range
	start := time.Now().AddDate(0, 0, -7)
	end := time.Now().AddDate(0, 0, 7)
	records, err := repo.GetByDateRange(start, end)
	if err != nil {
		log.Printf("Error getting attendance by date range: %v", err)
		return
	}
	fmt.Printf("✓ Found %d attendance records in date range\n", len(records))

	fmt.Println()
}

func testHomeworkRepository(repo homework.HomeworkRepository, courseRepo course.CourseRepository) {
	fmt.Println("📝 Testing Homework Repository...")

	courses, _ := courseRepo.GetAll(1, 0)
	if len(courses) == 0 {
		fmt.Println("⚠️  No courses found, skipping homework tests")
		return
	}

	// Create homework
	hw := &homework.Homework{
		Title:       "Calculus Problem Set 1",
		Description: "Complete problems 1-20",
		CourseID:    courses[0].ID,
		DueDate:     time.Now().AddDate(0, 0, 7),
		MaxScore:    100,
	}
	if err := repo.Create(hw); err != nil {
		log.Printf("Error creating homework: %v", err)
		return
	}
	fmt.Printf("✓ Created homework (ID: %d)\n", hw.ID)

	// Get upcoming
	upcoming, err := repo.GetUpcoming(10)
	if err != nil {
		log.Printf("Error getting upcoming homework: %v", err)
		return
	}
	fmt.Printf("✓ Found %d upcoming homework assignments\n", len(upcoming))

	fmt.Println()
}

func testExamRepository(repo exam.ExamRepository, courseRepo course.CourseRepository) {
	fmt.Println("📋 Testing Exam Repository...")

	courses, _ := courseRepo.GetAll(1, 0)
	if len(courses) == 0 {
		fmt.Println("⚠️  No courses found, skipping exam tests")
		return
	}

	// Create exam
	ex := &exam.Exam{
		Title:    "Midterm Exam",
		CourseID: courses[0].ID,
		ExamDate: time.Now().AddDate(0, 0, 14),
		Duration: 120,
		MaxScore: 100,
	}
	if err := repo.Create(ex); err != nil {
		log.Printf("Error creating exam: %v", err)
		return
	}
	fmt.Printf("✓ Created exam (ID: %d)\n", ex.ID)

	// Get upcoming
	upcoming, err := repo.GetUpcoming(10)
	if err != nil {
		log.Printf("Error getting upcoming exams: %v", err)
		return
	}
	fmt.Printf("✓ Found %d upcoming exams\n", len(upcoming))

	fmt.Println()
}

func testGradeRepository(repo grade.GradeRepository, studentRepo student.StudentRepository, examRepo exam.ExamRepository) {
	fmt.Println("🎓 Testing Grade Repository...")

	students, _ := studentRepo.GetAll(1, 0)
	exams, _ := examRepo.GetAll(1, 0)
	if len(students) == 0 || len(exams) == 0 {
		fmt.Println("⚠️  Missing dependencies, skipping grade tests")
		return
	}

	// Create grade
	g := &grade.Grade{
		StudentID: students[0].ID,
		ExamID:    exams[0].ID,
		Score:     95.5,
	}
	if err := repo.Create(g); err != nil {
		log.Printf("Error creating grade: %v", err)
		return
	}
	fmt.Printf("✓ Created grade (ID: %d, Score: %.1f)\n", g.ID, g.Score)

	// Get student average
	avg, err := repo.GetStudentAverage(students[0].ID)
	if err != nil {
		log.Printf("Error calculating average: %v", err)
		return
	}
	fmt.Printf("✓ Student average: %.2f\n", avg)

	fmt.Println()
}

func testSubmissionRepository(repo students_homework.StudentHomeworkRepository, studentRepo student.StudentRepository, homeworkRepo homework.HomeworkRepository) {
	fmt.Println("📤 Testing Submission Repository...")

	students, _ := studentRepo.GetAll(1, 0)
	homeworks, _ := homeworkRepo.GetAll(1, 0)
	// Check if submission already exists
	existing, err := repo.GetByStudentAndHomework(students[0].ID, homeworks[0].ID)
	if err == nil {
		fmt.Printf("⚠️  Submission already exists (ID: %d), skipping creation\n", existing.ID)
		return
	}
	if len(students) == 0 || len(homeworks) == 0 {
		fmt.Println("⚠️  Missing dependencies, skipping submission tests")
		return
	}

	// Create submission
	submissionTime := time.Now()
	score := 88.0
	submission := &students_homework.StudentHomework{
		StudentID:      students[0].ID,
		HomeworkID:     homeworks[0].ID,
		SubmissionDate: &submissionTime,
		Score:          &score,
		Status:         students_homework.HomeworkSubmitted,
	}
	if err := repo.Create(submission); err != nil {
		log.Printf("Error creating submission: %v", err)
		return
	}
	fmt.Printf("✓ Created submission (ID: %d, Score: %.1f)\n", submission.ID, *submission.Score)

	// Get pending submissions
	pending, err := repo.GetPendingByStudent(students[0].ID)
	if err != nil {
		log.Printf("Error getting pending submissions: %v", err)
		return
	}
	fmt.Printf("✓ Student has %d pending submissions\n", len(pending))

	fmt.Println()
}
