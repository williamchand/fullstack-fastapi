package jwt

import (
	"fmt"
	"log"
	"os"

	"github.com/williamchand/fullstack-fastapi/backend-go/config"
)

// NewService creates a JWT service from application config
func NewService(cfg *config.Config) (JWTService, error) {
	jwtConfig := cfg.GetJWTConfig()
	return NewServiceFromConfig(jwtConfig)
}

// NewServiceFromConfig creates a JWT service from JWT-specific config
func NewServiceFromConfig(cfg *config.JWTConfig) (JWTService, error) {
	// Try to create RSA signer first (recommended)
	if signer, err := createRSASigner(cfg); err == nil {
		log.Println("✅ Using RSA JWT signing (recommended for production)")
		return NewJWTService(
			signer,
			cfg.AccessTokenExpiration,
			cfg.RefreshTokenExpiration,
			cfg.Issuer,
		), nil
	}

	// Fallback to HMAC
	if cfg.HMACSecret != "" {
		log.Println("⚠️  Using HMAC JWT signing (consider using RSA for production)")
		signer := NewHMACSigner(cfg.HMACSecret)
		return NewJWTService(
			signer,
			cfg.AccessTokenExpiration,
			cfg.RefreshTokenExpiration,
			cfg.Issuer,
		), nil
	}

	return nil, fmt.Errorf("no JWT signing method available")
}

// createRSASigner attempts to create an RSA signer from various sources
func createRSASigner(cfg *config.JWTConfig) (Signer, error) {
	keyManager := NewKeyManager()

	// 1. Try loading from PEM strings
	if cfg.PrivateKeyPEM != "" && cfg.PublicKeyPEM != "" {
		if err := keyManager.LoadKeysFromStrings(cfg.PrivateKeyPEM, cfg.PublicKeyPEM); err != nil {
			return nil, fmt.Errorf("failed to load keys from environment: %w", err)
		}
	}

	// 2. Try loading from files
	if !keyManager.KeysExist() && cfg.PrivateKeyPath != "" && cfg.PublicKeyPath != "" {
		if fileExists(cfg.PrivateKeyPath) && fileExists(cfg.PublicKeyPath) {
			if err := keyManager.LoadKeysFromFiles(cfg.PrivateKeyPath, cfg.PublicKeyPath); err != nil {
				return nil, fmt.Errorf("failed to load keys from files: %w", err)
			}
		}
	}

	// 3. Auto-generate keys for development
	if !keyManager.KeysExist() && cfg.AutoGenerateKeys {
		log.Println("🔑 No JWT keys found, generating new RSA key pair for development...")
		if err := keyManager.GenerateKeys(); err != nil {
			return nil, fmt.Errorf("failed to generate keys: %w", err)
		}

		// Save generated keys
		if cfg.PrivateKeyPath != "" && cfg.PublicKeyPath != "" {
			if err := keyManager.SaveKeysToFiles(cfg.PrivateKeyPath, cfg.PublicKeyPath); err != nil {
				log.Printf("⚠️  Warning: Failed to save generated JWT keys: %v", err)
			} else {
				log.Printf("✅ Generated JWT keys saved to: %s, %s", cfg.PrivateKeyPath, cfg.PublicKeyPath)
			}
		}
	}

	if keyManager.KeysExist() {
		privateKey, publicKey := keyManager.GetKeys()
		return NewRSASigner(privateKey, publicKey), nil
	}

	return nil, fmt.Errorf("no RSA keys available")
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
