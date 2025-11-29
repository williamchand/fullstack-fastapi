package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type GuestRepository interface {
    TxProvider[GuestRepository]

    Add(ctx context.Context, g *entities.Guest) (*entities.Guest, error)
    Update(ctx context.Context, g *entities.Guest) (*entities.Guest, error)
    Delete(ctx context.Context, id uuid.UUID) error
    ListByWedding(ctx context.Context, weddingID uuid.UUID) ([]*entities.Guest, error)
}
