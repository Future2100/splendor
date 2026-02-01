package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string

	// JWT Configuration
	JWTSecret           string
	JWTAccessExpiry     int64 // seconds
	JWTRefreshExpiry    int64 // seconds

	// CORS Configuration
	AllowedOrigins []string

	// WebSocket Configuration
	WSReadTimeout  int64 // seconds
	WSWriteTimeout int64 // seconds
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/splendor?sslmode=disable"),

		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTAccessExpiry:  15 * 60,        // 15 minutes in seconds
		JWTRefreshExpiry: 7 * 24 * 60 * 60, // 7 days in seconds

		AllowedOrigins: []string{
			getEnv("FRONTEND_URL", "http://localhost:5173"),
		},

		WSReadTimeout:  60, // 60 seconds
		WSWriteTimeout: 10, // 10 seconds
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
