package student_courses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StudentCourseController handles HTTP requests for student course enrollments
type StudentCourseController struct {
	service StudentCourseService
}

// NewStudentCourseController creates a new student course controller
func NewStudentCourseController(service StudentCourseService) *StudentCourseController {
	return &StudentCourseController{service: service}
}

// Enroll enrolls a student in a course
func (c *StudentCourseController) Enroll(ctx *gin.Context) {
	var req EnrollStudentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Enroll(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetByID retrieves an enrollment by ID
func (c *StudentCourseController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	resp, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetByStudent retrieves enrollments for a student
func (c *StudentCourseController) GetByStudent(ctx *gin.Context) {
	studentID, err := strconv.ParseUint(ctx.Param("studentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	resp, err := c.service.GetByStudent(uint(studentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// GetByCourse retrieves enrollments for a course
func (c *StudentCourseController) GetByCourse(ctx *gin.Context) {
	courseID, err := strconv.ParseUint(ctx.Param("courseId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid course ID"})
		return
	}

	resp, err := c.service.GetByCourse(uint(courseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// Unenroll removes a student from a course
func (c *StudentCourseController) Unenroll(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Unenroll(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "student unenrolled successfully"})
}

// RegisterRoutes registers enrollment routes
func (c *StudentCourseController) RegisterRoutes(rg *gin.RouterGroup) {
	enrollments := rg.Group("/enrollments")
	{
		enrollments.POST("", c.Enroll)
		enrollments.GET("/:id", c.GetByID)
		enrollments.DELETE("/:id", c.Unenroll)
		enrollments.GET("/student/:studentId", c.GetByStudent)
		enrollments.GET("/course/:courseId", c.GetByCourse)
	}
}
