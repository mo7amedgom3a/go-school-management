package main

import (
	"fmt"
	"log"
	"net/http"

	"school_management/internal/config"

	"school_management/internal/database"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)

	// Create DB if missing
	database.CreateDatabaseIfNotExists(cfg)

	// Connect using GORM
	database.ConnectDB(cfg)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	// r := server.NewRouter()

	log.Println("🚀 Server starting on port:", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
