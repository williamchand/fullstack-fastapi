package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server Configuration
	HTTPPort string `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort string `envconfig:"GRPC_PORT" default:"9090"`
	Env      string `envconfig:"ENV" default:"development"`

	// Database Configuration
	DatabaseURL string `envconfig:"DB_URL" default:""`
	DBMaxConns  int    `envconfig:"DB_MAX_CONNS" default:"25"`
	DBMinConns  int    `envconfig:"DB_MIN_CONNS" default:"5"`

	// JWT Configuration - Multiple Options
	JWT struct {
		// RSA Key Files (Recommended for Production)
		PrivateKeyPath string `envconfig:"JWT_PRIVATE_KEY_PATH" default:"config/jwt/private.pem"`
		PublicKeyPath  string `envconfig:"JWT_PUBLIC_KEY_PATH" default:"config/jwt/public.pem"`

		// RSA Key Strings (Alternative for Docker/Cloud)
		PrivateKeyPEM string `envconfig:"JWT_PRIVATE_KEY"`
		PublicKeyPEM  string `envconfig:"JWT_PUBLIC_KEY"`

		// HMAC Fallback (Less Secure)
		HMACSecret string `envconfig:"JWT_HMAC_SECRET"`

		// Token Settings
		AccessTokenExpiration  time.Duration `envconfig:"JWT_ACCESS_EXPIRATION" default:"168h"`
		RefreshTokenExpiration time.Duration `envconfig:"JWT_REFRESH_EXPIRATION" default:"336h"` // 14 days
		Issuer                 string        `envconfig:"JWT_ISSUER" default:"salonapp"`

		// Auto-generation (Development Only)
		AutoGenerateKeys bool `envconfig:"JWT_AUTO_GENERATE_KEYS" default:"true"`
	}

	// OAuth Configuration
	OAuth struct {
		GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
		GoogleRedirectURL  string `envconfig:"GOOGLE_REDIRECT_URL" default:"http://localhost:8080/auth/google/callback"`
	}

	// Security Configuration
	Security struct {
		BCryptCost              int    `envconfig:"BCRYPT_COST" default:"10"`
		CORSAllowedOrigins      string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
		RateLimitRPS            int    `envconfig:"RATE_LIMIT_RPS" default:"100"`
		CredentialEncryptionKey string `envconfig:"CREDENTIAL_ENCRYPTION_KEY"`
	}

	// Logging Configuration
	Logging struct {
		Level    string `envconfig:"LOG_LEVEL" default:"info"`
		Format   string `envconfig:"LOG_FORMAT" default:"json"`
		FilePath string `envconfig:"LOG_FILE_PATH" default:""`
	}

	// Monitoring Configuration
	Monitoring struct {
		Enabled        bool   `envconfig:"MONITORING_ENABLED" default:"false"`
		PrometheusPort string `envconfig:"PROMETHEUS_PORT" default:"9091"`
	}

	// SMTP Configuration
	SMTP struct {
		Host     string `envconfig:"SMTP_HOST" default:""`
		Port     int    `envconfig:"SMTP_PORT" default:"587"`
		Username string `envconfig:"SMTP_USER" default:""`
		Password string `envconfig:"SMTP_PASSWORD" default:""`
		From     string `envconfig:"EMAILS_FROM_EMAIL" default:""`
	}

	// Superuser Configuration
	Superuser struct {
		Username  string `envconfig:"SUPERUSER_USERNAME"`
		Groupname string `envconfig:"SUPERUSER_GROUP" default:"superuser"`
		Password  string `envconfig:"SUPERUSER_PASSWORD"`
	}
    // Stripe Configuration
    Stripe struct {
        SecretKey       string `envconfig:"STRIPE_SECRET_KEY"`
        WebhookSecret   string `envconfig:"STRIPE_WEBHOOK_SECRET"`
        PriceID         string `envconfig:"STRIPE_PRICE_ID"`
    }

    // WAHA (WhatsApp Gateway) Configuration
    WAHA struct {
        URL     string `envconfig:"WAHA_URL" default:"http://localhost:3000"`
        APIKey  string `envconfig:"WAHA_API_KEY"`
        Session string `envconfig:"WAHA_SESSION" default:"default"`
    }

    // DOKU (Jokul Checkout) Configuration
    Doku struct {
        BaseURL  string `envconfig:"DOKU_BASE_URL" default:"https://api-sandbox.doku.com"`
        ClientID string `envconfig:"DOKU_CLIENT_ID"`
        SecretKey string `envconfig:"DOKU_SECRET_KEY"`
    }
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists (development)
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// UseRSAKeys returns true if RSA keys should be used for JWT
func (c *Config) UseRSAKeys() bool {
	return (c.JWT.PrivateKeyPath != "" && c.JWT.PublicKeyPath != "") ||
		(c.JWT.PrivateKeyPEM != "" && c.JWT.PublicKeyPEM != "")
}

// UseHMAC returns true if HMAC should be used for JWT
func (c *Config) UseHMAC() bool {
	return c.JWT.HMACSecret != "" && !c.UseRSAKeys()
}

// GetJWTConfig returns JWT-specific configuration
func (c *Config) GetJWTConfig() *JWTConfig {
	return &JWTConfig{
		PrivateKeyPath:         c.JWT.PrivateKeyPath,
		PublicKeyPath:          c.JWT.PublicKeyPath,
		PrivateKeyPEM:          c.JWT.PrivateKeyPEM,
		PublicKeyPEM:           c.JWT.PublicKeyPEM,
		HMACSecret:             c.JWT.HMACSecret,
		AccessTokenExpiration:  c.JWT.AccessTokenExpiration,
		RefreshTokenExpiration: c.JWT.RefreshTokenExpiration,
		Issuer:                 c.JWT.Issuer,
		AutoGenerateKeys:       c.JWT.AutoGenerateKeys,
	}
}

// GetJWTConfig returns JWT-specific configuration
func (c *Config) GetOauthConfig() *OAuthConfig {
	return &OAuthConfig{
		ClientID:     c.OAuth.GoogleClientID,
		ClientSecret: c.OAuth.GoogleClientSecret,
		RedirectURL:  c.OAuth.GoogleRedirectURL,
	}
}

// OAuthConfig is a subset of Config for OAuth-specific settings
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// JWTConfig is a subset of Config for JWT-specific settings
type JWTConfig struct {
	PrivateKeyPath         string
	PublicKeyPath          string
	PrivateKeyPEM          string
	PublicKeyPEM           string
	HMACSecret             string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	Issuer                 string
	AutoGenerateKeys       bool
}
