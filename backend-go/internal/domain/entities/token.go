package entities

import (
	"time"

	"github.com/google/uuid"
)

type TokenResult struct {
	Token     string
	ExpiresAt time.Time
}

// TokenClaims represents the claims embedded in JWT tokens
type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	ExpiresAt int64     `json:"exp"`
	IssuedAt  int64     `json:"iat"`
	Type      string    `json:"type"` // "access" or "refresh"
}
