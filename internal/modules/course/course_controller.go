package course

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CourseController handles HTTP requests for courses
type CourseController struct {
	service CourseService
}

// NewCourseController creates a new course controller
func NewCourseController(service CourseService) *CourseController {
	return &CourseController{service: service}
}

// Create creates a new course
func (c *CourseController) Create(ctx *gin.Context) {
	var req CreateCourseRequest

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

// GetByID retrieves a course by ID
func (c *CourseController) GetByID(ctx *gin.Context) {
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

// GetAll retrieves all courses with pagination
func (c *CourseController) GetAll(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	resp, err := c.service.GetAll(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   resp,
		"limit":  limit,
		"offset": offset,
		"count":  len(resp),
	})
}

// Update updates a course
func (c *CourseController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateCourseRequest
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

// Delete deletes a course
func (c *CourseController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "course deleted successfully"})
}

// GetByDepartment retrieves courses by department
func (c *CourseController) GetByDepartment(ctx *gin.Context) {
	deptID, err := strconv.ParseUint(ctx.Param("deptId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid department ID"})
		return
	}

	resp, err := c.service.GetByDepartment(uint(deptID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// RegisterRoutes registers course routes
func (c *CourseController) RegisterRoutes(rg *gin.RouterGroup) {
	courses := rg.Group("/courses")
	{
		courses.POST("", c.Create)
		courses.GET("/:id", c.GetByID)
		courses.GET("", c.GetAll)
		courses.PUT("/:id", c.Update)
		courses.DELETE("/:id", c.Delete)
		courses.GET("/department/:deptId", c.GetByDepartment)
	}
}
