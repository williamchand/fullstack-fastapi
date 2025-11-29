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
	IsEmailVerified bool
	IsPhoneVerified bool
	IsTOTPEnabled   bool
	TOTPSecret      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastLoginAt     *time.Time
	Roles           []string
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
    ProviderData   map[string]any
}

type DataSource struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    Name          string
    Type          string
    Host          string
    Port          int32
    DatabaseName  string
    Username      string
    PasswordEnc   string
    Options       map[string]any
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type AICredential struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Provider  string
    APIKeyEnc string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Subscription struct {
    ID                   uuid.UUID
    UserID               uuid.UUID
    StripeCustomerID     *string
    StripeSubscriptionID *string
    Status               string
    CurrentPeriodEnd     *time.Time
    CreatedAt            time.Time
    UpdatedAt            time.Time
}
