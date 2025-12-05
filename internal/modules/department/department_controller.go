package department

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DepartmentController handles HTTP requests for departments
type DepartmentController struct {
	service DepartmentService
}

// NewDepartmentController creates a new department controller
func NewDepartmentController(service DepartmentService) *DepartmentController {
	return &DepartmentController{service: service}
}

// Create creates a new department
func (c *DepartmentController) Create(ctx *gin.Context) {
	var req CreateDepartmentRequest

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

// GetByID retrieves a department by ID
func (c *DepartmentController) GetByID(ctx *gin.Context) {
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

// GetAll retrieves all departments with pagination
func (c *DepartmentController) GetAll(ctx *gin.Context) {
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

// Update updates a department
func (c *DepartmentController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req UpdateDepartmentRequest
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

// Delete deletes a department
func (c *DepartmentController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "department deleted successfully"})
}

// Search searches departments by name
func (c *DepartmentController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	resp, err := c.service.Search(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// RegisterRoutes registers department routes
func (c *DepartmentController) RegisterRoutes(rg *gin.RouterGroup) {
	departments := rg.Group("/departments")
	{
		departments.POST("", c.Create)
		departments.GET("/:id", c.GetByID)
		departments.GET("", c.GetAll)
		departments.PUT("/:id", c.Update)
		departments.DELETE("/:id", c.Delete)
		departments.GET("/search", c.Search)
	}
}
