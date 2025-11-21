package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// KeyManager handles RSA key generation, loading, and saving
type KeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewKeyManager creates a new key manager
func NewKeyManager() *KeyManager {
	return &KeyManager{}
}

// GenerateKeys generates a new RSA key pair
func (km *KeyManager) GenerateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	km.privateKey = privateKey
	km.publicKey = &privateKey.PublicKey
	return nil
}

// LoadKeysFromFiles loads keys from PEM files
func (km *KeyManager) LoadKeysFromFiles(privateKeyPath, publicKeyPath string) error {
	// Load private key
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := PEMToPrivateKey(privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load public key
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKey, err := PEMToPublicKey(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	km.privateKey = privateKey
	km.publicKey = publicKey
	return nil
}

// LoadKeysFromStrings loads keys from PEM strings
func (km *KeyManager) LoadKeysFromStrings(privateKeyPEM, publicKeyPEM string) error {
	privateKey, err := PEMToPrivateKey([]byte(privateKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := PEMToPublicKey([]byte(publicKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	km.privateKey = privateKey
	km.publicKey = publicKey
	return nil
}

// SaveKeysToFiles saves keys to PEM files
func (km *KeyManager) SaveKeysToFiles(privateKeyPath, publicKeyPath string) error {
	if km.privateKey == nil || km.publicKey == nil {
		return fmt.Errorf("no keys to save")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Save private key
	privateKeyPEM, err := PrivateKeyToPEM(km.privateKey)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	// Save public key
	publicKeyPEM, err := PublicKeyToPEM(km.publicKey)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %w", err)
	}

	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key file: %w", err)
	}

	return nil
}

// GetKeys returns the loaded keys
func (km *KeyManager) GetKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	return km.privateKey, km.publicKey
}

// KeysExist checks if keys are loaded
func (km *KeyManager) KeysExist() bool {
	return km.privateKey != nil && km.publicKey != nil
}

// PEM conversion utilities
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyPEM, nil
}

func PublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicKeyPEM, nil
}

func PEMToPrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func PEMToPublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPublicKey, nil
}
