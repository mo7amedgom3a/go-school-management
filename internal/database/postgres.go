package database

import (
	"database/sql"
	"fmt"
	"log"

	"school_management/internal/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL database:", cfg.DBName)
	return nil
}

func CreateDatabaseIfNotExists(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
	)

	// NOTE: pq driver is required here
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Cannot connect to PostgreSQL server: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + cfg.DBName)
	if err != nil {
		if err.Error() == fmt.Sprintf(`pq: database "%s" already exists`, cfg.DBName) {
			log.Println("ℹ️ Database already exists:", cfg.DBName)
			return nil
		}
		return fmt.Errorf("failed to create database: %v", err)
	}

	log.Println("🎉 Database created:", cfg.DBName)
	return nil
}
