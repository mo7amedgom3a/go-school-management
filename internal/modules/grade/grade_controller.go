package grade

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GradeController handles HTTP requests for grades
type GradeController struct {
	service GradeService
}

// NewGradeController creates a new grade controller
func NewGradeController(service GradeService) *GradeController {
	return &GradeController{service: service}
}

// Create creates a new grade
func (c *GradeController) Create(ctx *gin.Context) {
	var req CreateGradeRequest

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

// GetByID retrieves a grade by ID
func (c *GradeController) GetByID(ctx *gin.Context) {
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

// GetByStudent retrieves grades for a student
func (c *GradeController) GetByStudent(ctx *gin.Context) {
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

// GetByExam retrieves grades for an exam
func (c *GradeController) GetByExam(ctx *gin.Context) {
	examID, err := strconv.ParseUint(ctx.Param("examId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam ID"})
		return
	}

	resp, err := c.service.GetByExam(uint(examID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// GetStudentAverage calculates a student's average grade
func (c *GradeController) GetStudentAverage(ctx *gin.Context) {
	studentID, err := strconv.ParseUint(ctx.Param("studentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	avg, err := c.service.GetStudentAverage(uint(studentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"student_id": studentID, "average": avg})
}

// Update updates a grade
func (c *GradeController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateGradeRequest
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

// Delete deletes a grade
func (c *GradeController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "grade deleted successfully"})
}

// RegisterRoutes registers grade routes
func (c *GradeController) RegisterRoutes(rg *gin.RouterGroup) {
	grades := rg.Group("/grades")
	{
		grades.POST("", c.Create)
		grades.GET("/:id", c.GetByID)
		grades.PUT("/:id", c.Update)
		grades.DELETE("/:id", c.Delete)
		grades.GET("/student/:studentId", c.GetByStudent)
		grades.GET("/exam/:examId", c.GetByExam)
		grades.GET("/student/:studentId/average", c.GetStudentAverage)
	}
}
