package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort string `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort string `envconfig:"GRPC_PORT" default:"9090"`

	DatabaseURL string `envconfig:"DB_URL" default:"postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"`

	JWTSecret string        `envconfig:"JWT_SECRET" default:"your-default-secret-change-in-production"`
	JWTExpiry time.Duration `envconfig:"JWT_EXPIRY" default:"24h"`

	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `envconfig:"GOOGLE_REDIRECT_URL" default:"http://localhost:8080/auth/google/callback"`
}

func Load() (*Config, error) {
	// Load .env file automatically (if found)
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
