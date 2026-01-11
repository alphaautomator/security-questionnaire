package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	// Database configuration
	DatabaseURL string

	// AWS S3 configuration
	S3Bucket string
	S3Region string

	// AWS configuration
	AWSRegion string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3Region:    getEnvOrDefault("S3_REGION", "us-east-1"),
		AWSRegion:   getEnvOrDefault("AWS_REGION", "us-east-1"),
	}

	// Validate required configurations
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.S3Bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET environment variable is required")
	}

	return cfg, nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
