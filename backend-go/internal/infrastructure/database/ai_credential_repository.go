package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type aiCredentialRepository struct {
	queries *dbgen.Queries
	db      repositories.ConnectionPool
}

func NewAICredentialRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.AICredentialRepository {
	return &aiCredentialRepository{queries: q, db: db}
}

func (r *aiCredentialRepository) WithTx(tx pgx.Tx) repositories.AICredentialRepository {
	return &aiCredentialRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *aiCredentialRepository) Upsert(ctx context.Context, c *entities.AICredential) (*entities.AICredential, error) {
	out, err := r.queries.UpsertAICredential(ctx, dbgen.UpsertAICredentialParams{
		UserID:    c.UserID,
		Provider:  c.Provider,
		ApiKeyEnc: c.APIKeyEnc,
	})
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *aiCredentialRepository) Get(ctx context.Context, userID uuid.UUID, provider string) (*entities.AICredential, error) {
	out, err := r.queries.GetAICredential(ctx, dbgen.GetAICredentialParams{UserID: userID, Provider: provider})
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *aiCredentialRepository) toEntity(v *dbgen.AiCredential) *entities.AICredential {
	return &entities.AICredential{
		ID:        v.ID,
		UserID:    v.UserID,
		Provider:  v.Provider,
		APIKeyEnc: v.ApiKeyEnc,
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}
}
