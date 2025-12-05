package exam

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExamController handles HTTP requests for exams
type ExamController struct {
	service ExamService
}

// NewExamController creates a new exam controller
func NewExamController(service ExamService) *ExamController {
	return &ExamController{service: service}
}

// Create creates a new exam
func (c *ExamController) Create(ctx *gin.Context) {
	var req CreateExamRequest

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

// GetByID retrieves an exam by ID
func (c *ExamController) GetByID(ctx *gin.Context) {
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

// GetByCourse retrieves exams for a course
func (c *ExamController) GetByCourse(ctx *gin.Context) {
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

// GetUpcoming retrieves upcoming exams
func (c *ExamController) GetUpcoming(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.service.GetUpcoming(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// Update updates an exam
func (c *ExamController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateExamRequest
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

// Delete deletes an exam
func (c *ExamController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "exam deleted successfully"})
}

// RegisterRoutes registers exam routes
func (c *ExamController) RegisterRoutes(rg *gin.RouterGroup) {
	exams := rg.Group("/exams")
	{
		exams.POST("", c.Create)
		exams.GET("/:id", c.GetByID)
		exams.PUT("/:id", c.Update)
		exams.DELETE("/:id", c.Delete)
		exams.GET("/course/:courseId", c.GetByCourse)
		exams.GET("/upcoming", c.GetUpcoming)
	}
}
