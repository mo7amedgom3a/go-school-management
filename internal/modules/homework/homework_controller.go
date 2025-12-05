package homework

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HomeworkController handles HTTP requests for homework
type HomeworkController struct {
	service HomeworkService
}

// NewHomeworkController creates a new homework controller
func NewHomeworkController(service HomeworkService) *HomeworkController {
	return &HomeworkController{service: service}
}

// Create creates a new homework assignment
func (c *HomeworkController) Create(ctx *gin.Context) {
	var req CreateHomeworkRequest

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

// GetByID retrieves a homework assignment by ID
func (c *HomeworkController) GetByID(ctx *gin.Context) {
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

// GetByCourse retrieves homework assignments for a course
func (c *HomeworkController) GetByCourse(ctx *gin.Context) {
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

// GetUpcoming retrieves upcoming homework assignments
func (c *HomeworkController) GetUpcoming(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	resp, err := c.service.GetUpcoming(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// Update updates a homework assignment
func (c *HomeworkController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateHomeworkRequest
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

// Delete deletes a homework assignment
func (c *HomeworkController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "homework deleted successfully"})
}

// RegisterRoutes registers homework routes
func (c *HomeworkController) RegisterRoutes(rg *gin.RouterGroup) {
	homework := rg.Group("/homework")
	{
		homework.POST("", c.Create)
		homework.GET("/:id", c.GetByID)
		homework.PUT("/:id", c.Update)
		homework.DELETE("/:id", c.Delete)
		homework.GET("/course/:courseId", c.GetByCourse)
		homework.GET("/upcoming", c.GetUpcoming)
	}
}
