package entities

import (
    "time"

    "github.com/google/uuid"
)

type VerificationType string

const (
    VerificationTypeEmail VerificationType = "email"
    VerificationTypePhone VerificationType = "phone"
)

type VerificationCode struct {
    ID               uuid.UUID
    UserID           uuid.UUID
    Code             string
    Type             VerificationType
    ExtraMetadata    map[string]interface{}
    CreatedAt        time.Time
    ExpiresAt        time.Time
    UsedAt           *time.Time
}

type EmailTemplate struct {
    ID        uuid.UUID
    Name      string
    Subject   string
    Body      string
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}