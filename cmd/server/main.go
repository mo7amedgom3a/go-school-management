// @title School Management API
// @version 1.0
// @description REST API for School Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@schoolmanagement.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

package main

import (
	"log"

	"school_management/internal/config"
	"school_management/internal/database"
	"school_management/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create database if not exists
	database.CreateDatabaseIfNotExists(cfg)

	// Connect to database
	database.ConnectDB(cfg)

	// Run migrations
	database.RunMigrations()

	// Setup router
	router := server.SetupRouter()

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
