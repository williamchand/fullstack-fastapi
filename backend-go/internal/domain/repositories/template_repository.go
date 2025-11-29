package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type TemplateRepository interface {
    TxProvider[TemplateRepository]

    List(ctx context.Context) ([]*entities.Template, error)
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Template, error)
}
