package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
)

// jwtService implements JWTService
type jwtService struct {
	signer          Signer
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	issuer          string
}

// NewJWTService creates a new JWT service
func NewJWTService(signer Signer, accessTokenExp, refreshTokenExp time.Duration, issuer string) repositories.JWTRepository {
	return &jwtService{
		signer:          signer,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
		issuer:          issuer,
	}
}

// GenerateToken creates a new JWT access token
func (j *jwtService) GenerateToken(userID uuid.UUID, email string, roles []string) (*entities.TokenResult, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenExp)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"roles":   roles,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
		"iss":     j.issuer,
		"type":    "access",
	}

	token, err := j.signer.Sign(claims)
	return &entities.TokenResult{
		Token:     token,
		ExpiresAt: expiresAt,
	}, err
}

// GenerateRefreshToken creates a new refresh token
func (j *jwtService) GenerateRefreshToken(userID uuid.UUID) (*entities.TokenResult, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenExp)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
		"iss":     j.issuer,
		"type":    "refresh",
	}

	token, err := j.signer.Sign(claims)
	return &entities.TokenResult{
		Token:     token,
		ExpiresAt: expiresAt,
	}, err
}

// ValidateToken validates and parses a JWT token
func (j *jwtService) ValidateToken(tokenString string) (*entities.TokenClaims, error) {
	claims, err := j.signer.Verify(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// Extract and validate claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id claim is missing or invalid")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	email, _ := claims["email"].(string)

	var roles []string
	if rolesClaim, ok := claims["roles"]; ok {
		if rolesSlice, ok := rolesClaim.([]interface{}); ok {
			roles = make([]string, len(rolesSlice))
			for i, role := range rolesSlice {
				if roleStr, ok := role.(string); ok {
					roles[i] = roleStr
				}
			}
		}
	}

	exp, _ := claims["exp"].(float64)
	iat, _ := claims["iat"].(float64)
	tokenType, _ := claims["type"].(string)

	return &entities.TokenClaims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		ExpiresAt: int64(exp),
		IssuedAt:  int64(iat),
		Type:      tokenType,
	}, nil
}

// RefreshToken validates a refresh token and returns a new access token
func (j *jwtService) RefreshToken(refreshToken string) (*entities.TokenResult, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Generate new access token
	newAccessToken, err := j.GenerateToken(claims.UserID, claims.Email, claims.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}

// ExtractUserIDFromToken extracts user ID from token without full validation
func (j *jwtService) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	// Parse without validation to extract claims
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id claim is missing")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	return userID, nil
}
