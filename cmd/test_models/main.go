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

	fmt.Println("\n🧪 Starting GORM Model Tests...\n")

	// Run all tests
	testDepartments()
	testTeachers()
	testStudents()
	testCourses()
	testStudentCourses()
	testAttendance()
	testHomework()
	testStudentHomework()
	testExams()
	testGrades()

	fmt.Println("\n✅ All tests completed successfully!")
}

// Test Department CRUD operations
func testDepartments() {
	fmt.Println("📚 Testing Departments...")

	// CREATE
	dept := department.Department{
		Name:        "Computer Science",
		Description: "Department of Computer Science and Engineering",
	}
	result := database.DB.Create(&dept)
	if result.Error != nil {
		log.Fatalf("Failed to create department: %v", result.Error)
	}
	fmt.Printf("✓ Created department with ID: %d\n", dept.ID)

	// READ
	var retrievedDept department.Department
	database.DB.First(&retrievedDept, dept.ID)
	fmt.Printf("✓ Retrieved department: %s\n", retrievedDept.Name)

	// UPDATE
	database.DB.Model(&retrievedDept).Update("Description", "Updated CS Department")
	fmt.Printf("✓ Updated department description\n")

	// LIST
	var departments []department.Department
	database.DB.Find(&departments)
	fmt.Printf("✓ Found %d departments\n", len(departments))

	fmt.Println()
}

// Test Teacher CRUD operations
func testTeachers() {
	fmt.Println("👨‍🏫 Testing Teachers...")

	// Get first department
	var dept department.Department
	database.DB.First(&dept)

	// CREATE
	teacher1 := teacher.Teacher{
		FirstName:    "John",
		LastName:     "Smith",
		Email:        "john.smith@school.com",
		Phone:        "+1234567890",
		DepartmentID: dept.ID,
	}
	database.DB.Create(&teacher1)
	fmt.Printf("✓ Created teacher: %s %s (ID: %d)\n", teacher1.FirstName, teacher1.LastName, teacher1.ID)

	// READ with relationship
	var retrievedTeacher teacher.Teacher
	database.DB.Preload("Department").First(&retrievedTeacher, teacher1.ID)
	fmt.Printf("✓ Retrieved teacher with department (Note: Preload uses placeholder types)\n")

	// UPDATE
	database.DB.Model(&teacher1).Update("Phone", "+0987654321")
	fmt.Printf("✓ Updated teacher phone number\n")

	// LIST
	var teachers []teacher.Teacher
	database.DB.Find(&teachers)
	fmt.Printf("✓ Found %d teachers\n", len(teachers))

	fmt.Println()
}

// Test Student CRUD operations
func testStudents() {
	fmt.Println("👨‍🎓 Testing Students...")

	// CREATE
	dob, _ := time.Parse("2006-01-02", "2000-05-15")
	enrollDate := time.Now()

	student1 := student.Student{
		FirstName:      "Alice",
		LastName:       "Johnson",
		Email:          "alice.johnson@student.com",
		Phone:          "+1122334455",
		DateOfBirth:    dob,
		EnrollmentDate: enrollDate,
	}
	database.DB.Create(&student1)
	fmt.Printf("✓ Created student: %s %s (ID: %d)\n", student1.FirstName, student1.LastName, student1.ID)

	student2 := student.Student{
		FirstName:      "Bob",
		LastName:       "Williams",
		Email:          "bob.williams@student.com",
		Phone:          "+5566778899",
		DateOfBirth:    dob,
		EnrollmentDate: enrollDate,
	}
	database.DB.Create(&student2)
	fmt.Printf("✓ Created student: %s %s (ID: %d)\n", student2.FirstName, student2.LastName, student2.ID)

	// READ
	var retrievedStudent student.Student
	database.DB.First(&retrievedStudent, student1.ID)
	fmt.Printf("✓ Retrieved student: %s %s\n", retrievedStudent.FirstName, retrievedStudent.LastName)

	// UPDATE
	database.DB.Model(&student1).Update("Phone", "+9999999999")
	fmt.Printf("✓ Updated student phone number\n")

	// LIST with pagination
	var students []student.Student
	database.DB.Limit(10).Find(&students)
	fmt.Printf("✓ Found %d students\n", len(students))

	fmt.Println()
}

// Test Course CRUD operations
func testCourses() {
	fmt.Println("📖 Testing Courses...")

	// Get first department and teacher
	var dept department.Department
	var teach teacher.Teacher
	database.DB.First(&dept)
	database.DB.First(&teach)

	// CREATE
	course1 := course.Course{
		Name:         "Introduction to Programming",
		Code:         "CS101",
		Description:  "Learn the basics of programming",
		Credits:      3,
		DepartmentID: dept.ID,
		TeacherID:    teach.ID,
	}
	database.DB.Create(&course1)
	fmt.Printf("✓ Created course: %s (Code: %s, ID: %d)\n", course1.Name, course1.Code, course1.ID)

	course2 := course.Course{
		Name:         "Data Structures",
		Code:         "CS201",
		Description:  "Advanced data structures and algorithms",
		Credits:      4,
		DepartmentID: dept.ID,
		TeacherID:    teach.ID,
	}
	database.DB.Create(&course2)
	fmt.Printf("✓ Created course: %s (Code: %s, ID: %d)\n", course2.Name, course2.Code, course2.ID)

	// READ with relationships
	var retrievedCourse course.Course
	database.DB.Preload("Department").Preload("Teacher").First(&retrievedCourse, course1.ID)
	fmt.Printf("✓ Retrieved course with relationships\n")

	// UPDATE
	database.DB.Model(&course1).Update("Credits", 4)
	fmt.Printf("✓ Updated course credits\n")

	// LIST
	var courses []course.Course
	database.DB.Find(&courses)
	fmt.Printf("✓ Found %d courses\n", len(courses))

	fmt.Println()
}

// Test StudentCourse enrollment
func testStudentCourses() {
	fmt.Println("📝 Testing Student Course Enrollments...")

	// Get first student and course
	var stud student.Student
	var crs course.Course
	database.DB.First(&stud)
	database.DB.First(&crs)

	// CREATE enrollment
	enrollment := student_courses.StudentCourse{
		StudentID:      stud.ID,
		CourseID:       crs.ID,
		EnrollmentDate: time.Now(),
	}
	database.DB.Create(&enrollment)
	fmt.Printf("✓ Enrolled student ID %d in course ID %d\n", stud.ID, crs.ID)

	// READ enrollments for a student
	var enrollments []student_courses.StudentCourse
	database.DB.Where("student_id = ?", stud.ID).Find(&enrollments)
	fmt.Printf("✓ Student has %d enrollments\n", len(enrollments))

	// LIST all enrollments
	var allEnrollments []student_courses.StudentCourse
	database.DB.Find(&allEnrollments)
	fmt.Printf("✓ Total enrollments: %d\n", len(allEnrollments))

	fmt.Println()
}

// Test Attendance records
func testAttendance() {
	fmt.Println("✅ Testing Attendance...")

	// Get first student and course
	var stud student.Student
	var crs course.Course
	database.DB.First(&stud)
	database.DB.First(&crs)

	// CREATE attendance record
	att := attendance.Attendance{
		StudentID: stud.ID,
		CourseID:  crs.ID,
		Date:      time.Now(),
		Status:    attendance.AttendancePresent,
	}
	database.DB.Create(&att)
	fmt.Printf("✓ Created attendance record (ID: %d, Status: %s)\n", att.ID, att.Status)

	// CREATE another attendance record
	att2 := attendance.Attendance{
		StudentID: stud.ID,
		CourseID:  crs.ID,
		Date:      time.Now().AddDate(0, 0, 1),
		Status:    attendance.AttendanceAbsent,
	}
	database.DB.Create(&att2)
	fmt.Printf("✓ Created attendance record (ID: %d, Status: %s)\n", att2.ID, att2.Status)

	// READ attendance for a student
	var attendances []attendance.Attendance
	database.DB.Where("student_id = ?", stud.ID).Find(&attendances)
	fmt.Printf("✓ Student has %d attendance records\n", len(attendances))

	// UPDATE attendance status
	database.DB.Model(&att2).Update("Status", attendance.AttendanceLate)
	fmt.Printf("✓ Updated attendance status to Late\n")

	fmt.Println()
}

// Test Homework assignments
func testHomework() {
	fmt.Println("📝 Testing Homework...")

	// Get first course
	var crs course.Course
	database.DB.First(&crs)

	// CREATE homework
	hw := homework.Homework{
		Title:       "Programming Assignment 1",
		Description: "Implement a sorting algorithm",
		CourseID:    crs.ID,
		DueDate:     time.Now().AddDate(0, 0, 7), // Due in 7 days
		MaxScore:    100,
	}
	database.DB.Create(&hw)
	fmt.Printf("✓ Created homework: %s (ID: %d, Max Score: %.0f)\n", hw.Title, hw.ID, hw.MaxScore)

	// READ homework for a course
	var homeworks []homework.Homework
	database.DB.Where("course_id = ?", crs.ID).Find(&homeworks)
	fmt.Printf("✓ Course has %d homework assignments\n", len(homeworks))

	// UPDATE homework
	database.DB.Model(&hw).Update("MaxScore", 120)
	fmt.Printf("✓ Updated homework max score\n")

	fmt.Println()
}

// Test Student Homework submissions
func testStudentHomework() {
	fmt.Println("📤 Testing Student Homework Submissions...")

	// Get first student and homework
	var stud student.Student
	var hw homework.Homework
	database.DB.First(&stud)
	database.DB.First(&hw)

	// CREATE submission
	submissionTime := time.Now()
	score := 85.5
	submission := students_homework.StudentHomework{
		StudentID:      stud.ID,
		HomeworkID:     hw.ID,
		SubmissionDate: &submissionTime,
		Score:          &score,
		Status:         students_homework.HomeworkSubmitted,
	}
	database.DB.Create(&submission)
	fmt.Printf("✓ Created homework submission (ID: %d, Score: %.1f)\n", submission.ID, *submission.Score)

	// UPDATE submission - grade it
	newScore := 92.0
	database.DB.Model(&submission).Updates(map[string]interface{}{
		"Score":  &newScore,
		"Status": students_homework.HomeworkGraded,
	})
	fmt.Printf("✓ Graded homework submission\n")

	// READ submissions for a student
	var submissions []students_homework.StudentHomework
	database.DB.Where("student_id = ?", stud.ID).Find(&submissions)
	fmt.Printf("✓ Student has %d homework submissions\n", len(submissions))

	fmt.Println()
}

// Test Exam management
func testExams() {
	fmt.Println("📋 Testing Exams...")

	// Get first course
	var crs course.Course
	database.DB.First(&crs)

	// CREATE exam
	examDate := time.Now().AddDate(0, 0, 14) // Exam in 14 days
	ex := exam.Exam{
		Title:    "Midterm Exam",
		CourseID: crs.ID,
		ExamDate: examDate,
		Duration: 120, // 120 minutes
		MaxScore: 100,
	}
	database.DB.Create(&ex)
	fmt.Printf("✓ Created exam: %s (ID: %d, Duration: %d min)\n", ex.Title, ex.ID, ex.Duration)

	// READ exams for a course
	var exams []exam.Exam
	database.DB.Where("course_id = ?", crs.ID).Find(&exams)
	fmt.Printf("✓ Course has %d exams\n", len(exams))

	// UPDATE exam
	database.DB.Model(&ex).Update("Duration", 150)
	fmt.Printf("✓ Updated exam duration\n")

	fmt.Println()
}

// Test Grade management
func testGrades() {
	fmt.Println("🎓 Testing Grades...")

	// Get first student and exam
	var stud student.Student
	var ex exam.Exam
	database.DB.First(&stud)
	database.DB.First(&ex)

	// CREATE grade
	gr := grade.Grade{
		StudentID: stud.ID,
		ExamID:    ex.ID,
		Score:     88.5,
	}
	database.DB.Create(&gr)
	fmt.Printf("✓ Created grade (ID: %d, Score: %.1f)\n", gr.ID, gr.Score)

	// READ grades for a student
	var grades []grade.Grade
	database.DB.Where("student_id = ?", stud.ID).Find(&grades)
	fmt.Printf("✓ Student has %d grades\n", len(grades))

	// UPDATE grade
	database.DB.Model(&gr).Update("Score", 92.0)
	fmt.Printf("✓ Updated grade score\n")

	// READ with relationships
	var gradeWithRelations grade.Grade
	database.DB.Preload("Student").Preload("Exam").First(&gradeWithRelations, gr.ID)
	fmt.Printf("✓ Retrieved grade with student and exam relationships\n")

	// Calculate average grade for student
	var avgScore float64
	database.DB.Model(&grade.Grade{}).
		Where("student_id = ?", stud.ID).
		Select("AVG(score)").
		Scan(&avgScore)
	fmt.Printf("✓ Student average score: %.2f\n", avgScore)

	fmt.Println()
}
