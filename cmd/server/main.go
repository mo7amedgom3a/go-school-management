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
