package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type PaymentRepository interface {
	TxProvider[PaymentRepository]

	Create(ctx context.Context, p *entities.Payment) (*entities.Payment, error)
	GetByTransaction(ctx context.Context, txid string) (*entities.Payment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)
	UpdateStatus(ctx context.Context, txid string, status entities.PaymentStatus, amount *float64, currency *string, metadata map[string]any) (*entities.Payment, error)
}
