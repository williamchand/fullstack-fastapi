package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type weddingRepository struct {
	queries *dbgen.Queries
	db      repositories.ConnectionPool
}

func NewWeddingRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.WeddingRepository {
	return &weddingRepository{queries: q, db: db}
}

func (r *weddingRepository) WithTx(tx pgx.Tx) repositories.WeddingRepository {
	return &weddingRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *weddingRepository) Create(ctx context.Context, w *entities.Wedding) (*entities.Wedding, error) {
	out, err := r.queries.CreateWedding(ctx, dbgen.CreateWeddingParams{
		UserID:       w.UserID,
		TemplateID:   toPgUUIDPtr(w.TemplateID),
		PaymentID:    toPgUUIDPtr(w.PaymentID),
		Status:       dbgen.WeddingStatus(w.Status),
		CustomDomain: toPgText(w.CustomDomain),
		Slug:         toPgText(w.Slug),
		ConfigData:   toPgJSON(w.ConfigData),
	})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Wedding, error) {
	out, err := r.queries.GetWeddingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Wedding, error) {
	rows, err := r.queries.GetWeddingsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*entities.Wedding, 0, len(rows))
	for _, v := range rows {
		res = append(res, r.toEntity(v))
	}
	return res, nil
}

func (r *weddingRepository) UpdateConfig(ctx context.Context, id uuid.UUID, config map[string]any) (*entities.Wedding, error) {
	out, err := r.queries.UpdateWeddingConfig(ctx, dbgen.UpdateWeddingConfigParams{ID: id, ConfigData: toPgJSON(config)})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) SetTemplate(ctx context.Context, id uuid.UUID, templateID uuid.UUID) (*entities.Wedding, error) {
	out, err := r.queries.SetWeddingTemplate(ctx, dbgen.SetWeddingTemplateParams{ID: id, TemplateID: toPgUUIDPtr(&templateID)})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) SetPayment(ctx context.Context, id uuid.UUID, paymentID uuid.UUID) (*entities.Wedding, error) {
	out, err := r.queries.SetWeddingPayment(ctx, dbgen.SetWeddingPaymentParams{ID: id, PaymentID: toPgUUIDPtr(&paymentID)})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) SetDomain(ctx context.Context, id uuid.UUID, domain string) (*entities.Wedding, error) {
	out, err := r.queries.SetWeddingDomain(ctx, dbgen.SetWeddingDomainParams{ID: id, CustomDomain: toPgText(&domain)})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) SetSlug(ctx context.Context, id uuid.UUID, slug string) (*entities.Wedding, error) {
	out, err := r.queries.SetWeddingSlug(ctx, dbgen.SetWeddingSlugParams{ID: id, Slug: toPgText(&slug)})
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) Publish(ctx context.Context, id uuid.UUID) (*entities.Wedding, error) {
	out, err := r.queries.PublishWedding(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toEntity(out), nil
}

func (r *weddingRepository) toEntity(v dbgen.Wedding) *entities.Wedding {
	var tmplID *uuid.UUID
	if v.TemplateID.Valid {
		id := uuid.UUID(v.TemplateID.Bytes)
		tmplID = &id
	}
	var payID *uuid.UUID
	if v.PaymentID.Valid {
		id := uuid.UUID(v.PaymentID.Bytes)
		payID = &id
	}
	return &entities.Wedding{
		ID:           v.ID,
		UserID:       v.UserID,
		TemplateID:   tmplID,
		PaymentID:    payID,
		Status:       entities.WeddingStatus(v.Status),
		CustomDomain: fromPgText(v.CustomDomain),
		Slug:         fromPgText(v.Slug),
		ConfigData:   fromPgJSON(v.ConfigData),
		CreatedAt:    v.CreatedAt.Time,
		DeletedAt:    fromPgTime(v.DeletedAt),
	}
}
