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

// Create godoc
// @Summary Create department
// @Description Create a new department
// @Tags departments
// @Accept json
// @Produce json
// @Param department body CreateDepartmentRequest true "Department data"
// @Success 201 {object} DepartmentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /departments [post]
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

// GetByID godoc
// @Summary Get department by ID
// @Description Get a single department by its ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} DepartmentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /departments/{id} [get]
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

// GetAll godoc
// @Summary List departments
// @Description Get all departments with pagination
// @Tags departments
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /departments [get]
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

// Update godoc
// @Summary Update department
// @Description Update an existing department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param department body UpdateDepartmentRequest true "Department data"
// @Success 200 {object} DepartmentResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /departments/{id} [put]
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

// Delete godoc
// @Summary Delete department
// @Description Delete a department by ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /departments/{id} [delete]
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

// Search godoc
// @Summary Search departments
// @Description Search departments by name
// @Tags departments
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /departments/search [get]
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
