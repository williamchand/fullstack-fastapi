package entities

import (
    "time"
    "github.com/google/uuid"
)

type RSVPStatus string

const (
    RSVPYes   RSVPStatus = "yes"
    RSVPNo    RSVPStatus = "no"
    RSVPMaybe RSVPStatus = "maybe"
)

type WeddingStatus string

const (
    WeddingDraft   WeddingStatus = "draft"
    WeddingActive  WeddingStatus = "active"
    WeddingArchived WeddingStatus = "archived"
)

type Template struct {
    ID           uuid.UUID
    Name         string
    ThemeConfig  map[string]any
    ConfigSchema map[string]any
    PreviewURL   *string
    Price        float64
    CreatedAt    time.Time
}

type Wedding struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    TemplateID    *uuid.UUID
    PaymentID     *uuid.UUID
    Status        WeddingStatus
    CustomDomain  *string
    Slug          *string
    ConfigData    map[string]any
    CreatedAt     time.Time
    DeletedAt     *time.Time
}

type Guest struct {
    ID         uuid.UUID
    WeddingID  uuid.UUID
    Name       string
    Contact    string
    RSVPStatus RSVPStatus
    Message    *string
    CreatedAt  time.Time
    DeletedAt  *time.Time
}
