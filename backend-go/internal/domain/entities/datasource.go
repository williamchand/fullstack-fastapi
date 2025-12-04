package entities

import (
	"time"

	"github.com/google/uuid"
)

type DataSource struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	Type         string
	Host         string
	Port         int32
	DatabaseName string
	Username     string
	PasswordEnc  string
	Options      map[string]any
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AICredential struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Provider  string
	APIKeyEnc string
	CreatedAt time.Time
	UpdatedAt time.Time
}
