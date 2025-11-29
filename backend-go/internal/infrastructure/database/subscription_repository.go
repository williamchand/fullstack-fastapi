package database

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type subscriptionRepository struct {
    queries *dbgen.Queries
    db      repositories.ConnectionPool
}

func NewSubscriptionRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.SubscriptionRepository {
    return &subscriptionRepository{queries: q, db: db}
}

func (r *subscriptionRepository) WithTx(tx pgx.Tx) repositories.SubscriptionRepository {
    return &subscriptionRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *subscriptionRepository) Upsert(ctx context.Context, s *entities.Subscription) (*entities.Subscription, error) {
    out, err := r.queries.UpsertSubscription(ctx, dbgen.UpsertSubscriptionParams{
        UserID:             s.UserID,
        StripeCustomerID:   toPgText(s.StripeCustomerID),
        StripeSubscriptionID: toPgText(s.StripeSubscriptionID),
        Status:             s.Status,
        CurrentPeriodEnd:   toPgTimestamptz(s.CurrentPeriodEnd),
    })
    if err != nil {
        return nil, err
    }
    return r.toEntity(&out), nil
}

func (r *subscriptionRepository) GetByUser(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error) {
    out, err := r.queries.GetSubscriptionByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    return r.toEntity(&out), nil
}

func (r *subscriptionRepository) toEntity(v *dbgen.Subscription) *entities.Subscription {
    return &entities.Subscription{
        ID:                   v.ID,
        UserID:               v.UserID,
        StripeCustomerID:     fromPgText(v.StripeCustomerID),
        StripeSubscriptionID: fromPgText(v.StripeSubscriptionID),
        Status:               v.Status,
        CurrentPeriodEnd:     fromPgTime(v.CurrentPeriodEnd),
        CreatedAt:            v.CreatedAt.Time,
        UpdatedAt:            v.UpdatedAt.Time,
    }
}
