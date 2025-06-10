package db

import (
	"fmt"
	"time"

	"github.com/uvish/url-shortener/internal/config"
	"github.com/uvish/url-shortener/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(config config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode)

	maxRetries := 10
	retryDelay := 5 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			if err := db.AutoMigrate(&model.URL{}); err != nil {
				return fmt.Errorf("failed to run auto-migration: %w", err)
			}
			DB = db
			return nil
		}
		if attempt == maxRetries {
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}
		fmt.Printf("Database connection attempt %d failed: %v. Retrying in %v...\n", attempt, err, retryDelay)
		time.Sleep(retryDelay)
	}
	return nil
}
