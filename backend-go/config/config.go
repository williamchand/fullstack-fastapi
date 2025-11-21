package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	HTTPPort string
	GRPCPort string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret string
	JWTExpiry time.Duration

	// OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		GRPCPort: getEnv("GRPC_PORT", "9090"),

		DatabaseURL: getEnv("DB_URL", "postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"),

		JWTSecret: getEnv("JWT_SECRET", "your-default-secret-change-in-production"),
		JWTExpiry: getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
