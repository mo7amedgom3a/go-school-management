package students_homework

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StudentHomeworkController handles HTTP requests for student homework submissions
type StudentHomeworkController struct {
	service StudentHomeworkService
}

// NewStudentHomeworkController creates a new student homework controller
func NewStudentHomeworkController(service StudentHomeworkService) *StudentHomeworkController {
	return &StudentHomeworkController{service: service}
}

// Submit submits homework
func (c *StudentHomeworkController) Submit(ctx *gin.Context) {
	var req SubmitHomeworkRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Submit(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// Grade grades a homework submission
func (c *StudentHomeworkController) Grade(ctx *gin.Context) {
	var req GradeHomeworkRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.service.Grade(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetByID retrieves a submission by ID
func (c *StudentHomeworkController) GetByID(ctx *gin.Context) {
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

// GetByStudent retrieves submissions for a student
func (c *StudentHomeworkController) GetByStudent(ctx *gin.Context) {
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

// GetByHomework retrieves submissions for a homework
func (c *StudentHomeworkController) GetByHomework(ctx *gin.Context) {
	homeworkID, err := strconv.ParseUint(ctx.Param("homeworkId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid homework ID"})
		return
	}

	resp, err := c.service.GetByHomework(uint(homeworkID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// GetPendingByStudent retrieves pending submissions for a student
func (c *StudentHomeworkController) GetPendingByStudent(ctx *gin.Context) {
	studentID, err := strconv.ParseUint(ctx.Param("studentId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	resp, err := c.service.GetPendingByStudent(uint(studentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": resp, "count": len(resp)})
}

// Delete deletes a submission
func (c *StudentHomeworkController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "submission deleted successfully"})
}

// RegisterRoutes registers submission routes
func (c *StudentHomeworkController) RegisterRoutes(rg *gin.RouterGroup) {
	submissions := rg.Group("/submissions")
	{
		submissions.POST("", c.Submit)
		submissions.POST("/grade", c.Grade)
		submissions.GET("/:id", c.GetByID)
		submissions.DELETE("/:id", c.Delete)
		submissions.GET("/student/:studentId", c.GetByStudent)
		submissions.GET("/homework/:homeworkId", c.GetByHomework)
		submissions.GET("/student/:studentId/pending", c.GetPendingByStudent)
	}
}
