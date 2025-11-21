package jwt

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Signer defines the interface for JWT signing and verification
type Signer interface {
	Sign(claims jwt.MapClaims) (string, error)
	Verify(tokenString string) (jwt.MapClaims, error)
}

// RSASigner uses RSA keys for signing
type RSASigner struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSASigner creates a new RSA signer
func NewRSASigner(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) Signer {
	return &RSASigner{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (r *RSASigner) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(r.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token with RSA: %w", err)
	}
	return signedToken, nil
}

func (r *RSASigner) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return r.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// HMACSigner uses HMAC for signing
type HMACSigner struct {
	secret []byte
}

// NewHMACSigner creates a new HMAC signer
func NewHMACSigner(secret string) Signer {
	return &HMACSigner{
		secret: []byte(secret),
	}
}

func (h *HMACSigner) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(h.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token with HMAC: %w", err)
	}
	return signedToken, nil
}

func (h *HMACSigner) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
