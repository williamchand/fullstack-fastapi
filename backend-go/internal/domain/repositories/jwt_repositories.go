package repositories

import (
	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

// JWTRepository defines the interface for JWT token operations
type JWTRepository interface {
	GenerateToken(userID uuid.UUID, email string, roles []string) (*entities.TokenResult, error)
	GenerateRefreshToken(userID uuid.UUID) (*entities.TokenResult, error)
	ValidateToken(tokenString string) (*entities.TokenClaims, error)
	RefreshToken(refreshToken string) (*entities.TokenResult, error)
	ExtractUserIDFromToken(tokenString string) (uuid.UUID, error)
}
