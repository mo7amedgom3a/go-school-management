package server

import (
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

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "school_management/docs" // Import generated docs
)

// SetupRouter creates and configures the Gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/api/v1")

	// Initialize repositories
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

	// Initialize services
	deptService := department.NewDepartmentService(deptRepo)
	teacherService := teacher.NewTeacherService(teacherRepo)
	studentService := student.NewStudentService(studentRepo)
	courseService := course.NewCourseService(courseRepo)
	attendanceService := attendance.NewAttendanceService(attendanceRepo)
	homeworkService := homework.NewHomeworkService(homeworkRepo)
	examService := exam.NewExamService(examRepo)
	gradeService := grade.NewGradeService(gradeRepo)
	enrollmentService := student_courses.NewStudentCourseService(enrollmentRepo)
	submissionService := students_homework.NewStudentHomeworkService(submissionRepo)

	// Initialize controllers
	deptController := department.NewDepartmentController(deptService)
	teacherController := teacher.NewTeacherController(teacherService)
	studentController := student.NewStudentController(studentService)
	courseController := course.NewCourseController(courseService)
	attendanceController := attendance.NewAttendanceController(attendanceService)
	homeworkController := homework.NewHomeworkController(homeworkService)
	examController := exam.NewExamController(examService)
	gradeController := grade.NewGradeController(gradeService)
	enrollmentController := student_courses.NewStudentCourseController(enrollmentService)
	submissionController := students_homework.NewStudentHomeworkController(submissionService)

	// Register routes
	deptController.RegisterRoutes(v1)
	teacherController.RegisterRoutes(v1)
	studentController.RegisterRoutes(v1)
	courseController.RegisterRoutes(v1)
	attendanceController.RegisterRoutes(v1)
	homeworkController.RegisterRoutes(v1)
	examController.RegisterRoutes(v1)
	gradeController.RegisterRoutes(v1)
	enrollmentController.RegisterRoutes(v1)
	submissionController.RegisterRoutes(v1)

	return router
}
