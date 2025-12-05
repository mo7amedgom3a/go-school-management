package teacher

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TeacherController handles HTTP requests for teachers
type TeacherController struct {
	service TeacherService
}

// NewTeacherController creates a new teacher controller
func NewTeacherController(service TeacherService) *TeacherController {
	return &TeacherController{service: service}
}

// Create creates a new teacher
func (c *TeacherController) Create(ctx *gin.Context) {
	var req CreateTeacherRequest

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

// GetByID retrieves a teacher by ID
func (c *TeacherController) GetByID(ctx *gin.Context) {
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

// GetAll retrieves all teachers with pagination
func (c *TeacherController) GetAll(ctx *gin.Context) {
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

// Update updates a teacher
func (c *TeacherController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateTeacherRequest
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

// Delete deletes a teacher
func (c *TeacherController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "teacher deleted successfully"})
}

// GetByDepartment retrieves teachers by department
func (c *TeacherController) GetByDepartment(ctx *gin.Context) {
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

// RegisterRoutes registers teacher routes
func (c *TeacherController) RegisterRoutes(rg *gin.RouterGroup) {
	teachers := rg.Group("/teachers")
	{
		teachers.POST("", c.Create)
		teachers.GET("/:id", c.GetByID)
		teachers.GET("", c.GetAll)
		teachers.PUT("/:id", c.Update)
		teachers.DELETE("/:id", c.Delete)
		teachers.GET("/department/:deptId", c.GetByDepartment)
	}
}
