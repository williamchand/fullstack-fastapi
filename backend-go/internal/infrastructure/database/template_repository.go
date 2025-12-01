package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type templateRepository struct {
	queries *dbgen.Queries
	db      repositories.ConnectionPool
}

func NewTemplateRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.TemplateRepository {
	return &templateRepository{queries: q, db: db}
}

func (r *templateRepository) WithTx(tx pgx.Tx) repositories.TemplateRepository {
	return &templateRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *templateRepository) List(ctx context.Context) ([]*entities.Template, error) {
	rows, err := r.queries.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*entities.Template, 0, len(rows))
	for _, v := range rows {
		res = append(res, r.toEntity(v))
	}
	return res, nil
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Template, error) {
	v, err := r.queries.GetTemplateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toEntity(v), nil
}

func (r *templateRepository) toEntity(v dbgen.Template) *entities.Template {
	return &entities.Template{
		ID:           v.ID,
		Name:         v.Name,
		ThemeConfig:  fromPgJSON(v.ThemeConfig),
		ConfigSchema: fromPgJSON(v.ConfigSchema),
		PreviewURL:   fromPgText(v.PreviewUrl),
		Price:        fromPgNumericToFloat64(v.Price),
		CreatedAt:    v.CreatedAt.Time,
	}
}
