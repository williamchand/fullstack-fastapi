package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type WeddingRepository interface {
    TxProvider[WeddingRepository]

    Create(ctx context.Context, w *entities.Wedding) (*entities.Wedding, error)
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Wedding, error)
    ListByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Wedding, error)
    UpdateConfig(ctx context.Context, id uuid.UUID, config map[string]any) (*entities.Wedding, error)
    SetTemplate(ctx context.Context, id uuid.UUID, templateID uuid.UUID) (*entities.Wedding, error)
    SetPayment(ctx context.Context, id uuid.UUID, paymentID uuid.UUID) (*entities.Wedding, error)
    SetDomain(ctx context.Context, id uuid.UUID, domain string) (*entities.Wedding, error)
    SetSlug(ctx context.Context, id uuid.UUID, slug string) (*entities.Wedding, error)
    Publish(ctx context.Context, id uuid.UUID) (*entities.Wedding, error)
}
