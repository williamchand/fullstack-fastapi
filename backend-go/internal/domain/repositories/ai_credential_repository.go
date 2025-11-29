package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type AICredentialRepository interface {
    TxProvider[AICredentialRepository]

    Upsert(ctx context.Context, c *entities.AICredential) (*entities.AICredential, error)
    Get(ctx context.Context, userID uuid.UUID, provider string) (*entities.AICredential, error)
}
