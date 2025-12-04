package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
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

type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
	PaymentProviderDoku   PaymentProvider = "doku"
)
