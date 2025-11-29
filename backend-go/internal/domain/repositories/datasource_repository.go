package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type DataSourceRepository interface {
    TxProvider[DataSourceRepository]

    Create(ctx context.Context, ds *entities.DataSource) (*entities.DataSource, error)
    GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DataSource, error)
}
