package database

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type emailTemplateRepository struct {
    queries *dbgen.Queries
    db      repositories.ConnectionPool
}

func NewEmailTemplateRepository(queries *dbgen.Queries, db repositories.ConnectionPool) repositories.EmailTemplateRepository {
    return &emailTemplateRepository{queries: queries, db: db}
}

func (r *emailTemplateRepository) WithTx(tx pgx.Tx) repositories.EmailTemplateRepository {
    return &emailTemplateRepository{
        queries: r.queries.WithTx(tx),
        db:      r.db,
    }
}

func (r *emailTemplateRepository) GetByName(ctx context.Context, name string) (*entities.EmailTemplate, error) {
    tpl, err := r.queries.GetEmailTemplateByName(ctx, name)
    if err != nil {
        return nil, err
    }
    return &entities.EmailTemplate{
        ID:        tpl.ID,
        Name:      tpl.Name,
        Subject:   tpl.Subject,
        Body:      tpl.Body,
        IsActive:  tpl.IsActive,
        CreatedAt: tpl.CreatedAt.Time,
        UpdatedAt: tpl.UpdatedAt.Time,
    }, nil
}