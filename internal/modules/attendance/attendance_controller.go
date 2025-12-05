package attendance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AttendanceController handles HTTP requests for attendance
type AttendanceController struct {
	service AttendanceService
}

// NewAttendanceController creates a new attendance controller
func NewAttendanceController(service AttendanceService) *AttendanceController {
	return &AttendanceController{service: service}
}

// Create creates a new attendance record
func (c *AttendanceController) Create(ctx *gin.Context) {
	var req CreateAttendanceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetByID retrieves an attendance record by ID
func (c *AttendanceController) GetByID(ctx *gin.Context) {
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

// GetByStudent retrieves attendance records for a student
func (c *AttendanceController) GetByStudent(ctx *gin.Context) {
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

// GetByCourse retrieves attendance records for a course
func (c *AttendanceController) GetByCourse(ctx *gin.Context) {
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

// Update updates an attendance record
func (c *AttendanceController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Update(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// Delete deletes an attendance record
func (c *AttendanceController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "attendance deleted successfully"})
}

// RegisterRoutes registers attendance routes
func (c *AttendanceController) RegisterRoutes(rg *gin.RouterGroup) {
	attendance := rg.Group("/attendance")
	{
		attendance.POST("", c.Create)
		attendance.GET("/:id", c.GetByID)
		attendance.PUT("/:id", c.Update)
		attendance.DELETE("/:id", c.Delete)
		attendance.GET("/student/:studentId", c.GetByStudent)
		attendance.GET("/course/:courseId", c.GetByCourse)
	}
}
