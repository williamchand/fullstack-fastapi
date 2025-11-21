package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
)

func main() {
	// Create key manager
	keyManager := jwt.NewKeyManager()

	// Generate keys
	if err := keyManager.GenerateKeys(); err != nil {
		log.Fatal("Failed to generate keys:", err)
	}

	// Create config directory
	configDir := "config/jwt"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Failed to create config directory:", err)
	}

	// Save keys
	privateKeyPath := filepath.Join(configDir, "private.pem")
	publicKeyPath := filepath.Join(configDir, "public.pem")

	if err := keyManager.SaveKeysToFiles(privateKeyPath, publicKeyPath); err != nil {
		log.Fatal("Failed to save keys:", err)
	}

	fmt.Printf("JWT keys generated successfully!\n")
	fmt.Printf("Private key: %s\n", privateKeyPath)
	fmt.Printf("Public key:  %s\n", publicKeyPath)
	fmt.Printf("\nAdd these to your .env file:\n")
	fmt.Printf("JWT_PRIVATE_KEY_PATH=%s\n", privateKeyPath)
	fmt.Printf("JWT_PUBLIC_KEY_PATH=%s\n", publicKeyPath)
}
