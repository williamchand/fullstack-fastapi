package repositories

import (
    "context"
    "github.com/google/uuid"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type SubscriptionRepository interface {
    TxProvider[SubscriptionRepository]

    Upsert(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error)
    GetByUser(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error)
}
