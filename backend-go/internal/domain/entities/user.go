package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Email           string
	PhoneNumber     *string
	FullName        *string
	HashedPassword  *string
	IsActive        bool
	IsSuperuser     bool
	IsEmailVerified bool
	IsPhoneVerified bool
	IsTOTPEnabled   bool
	TOTPSecret      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastLoginAt     *time.Time
	Roles           []Role
}

type Role struct {
	ID          int32
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OAuthAccount struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	AccessToken    *string
	RefreshToken   *string
	TokenExpiresAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProviderData   map[string]interface{}
}
