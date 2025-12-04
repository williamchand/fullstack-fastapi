package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusActive    PaymentStatus = "active"
)

type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderDoku   PaymentProvider = "doku"
)

type Payment struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	PaymentMethodID *uuid.UUID
	Provider        PaymentProvider
	Amount          float64
	Currency        string
	Status          PaymentStatus
	TransactionID   string
	ExtraMetadata   map[string]any
	CreatedAt       time.Time
}

type Subscription struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	StripeCustomerID     *string
	StripeSubscriptionID *string
	Status               PaymentStatus
	CurrentPeriodEnd     *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
