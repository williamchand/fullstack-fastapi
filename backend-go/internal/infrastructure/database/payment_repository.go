package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type paymentRepository struct {
	queries *dbgen.Queries
	db      repositories.ConnectionPool
}

func NewPaymentRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.PaymentRepository {
	return &paymentRepository{queries: q, db: db}
}

func (r *paymentRepository) WithTx(tx pgx.Tx) repositories.PaymentRepository {
	return &paymentRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *paymentRepository) Create(ctx context.Context, p *entities.Payment) (*entities.Payment, error) {
	out, err := r.queries.CreatePayment(ctx, dbgen.CreatePaymentParams{
		UserID:          p.UserID,
		PaymentMethodID: toPgUUIDPtr(p.PaymentMethodID),
		Amount:          p.Amount,
		Currency:        p.Currency,
		Status:          string(p.Status),
		TransactionID:   p.TransactionID,
		ExtraMetadata:   toPgJSON(p.ExtraMetadata),
	})
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *paymentRepository) GetByTransaction(ctx context.Context, txid string) (*entities.Payment, error) {
	out, err := r.queries.GetPaymentByTransaction(ctx, txid)
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	out, err := r.queries.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, txid string, status entities.PaymentStatus, amount *float64, currency *string, metadata map[string]any) (*entities.Payment, error) {
	out, err := r.queries.UpdatePaymentStatus(ctx, dbgen.UpdatePaymentStatusParams{
		TransactionID: txid,
		Status:        string(status),
		Amount:        toPgNumeric(amount),
		Currency:      toPgText(currency),
		ExtraMetadata: toPgJSON(metadata),
	})
	if err != nil {
		return nil, err
	}
	return r.toEntity(&out), nil
}

func (r *paymentRepository) toEntity(v *dbgen.Payment) *entities.Payment {
	var pmID *uuid.UUID
	if v.PaymentMethodID.Valid {
		id := uuid.UUID(v.PaymentMethodID.Bytes)
		pmID = &id
	}
	return &entities.Payment{
		ID:              v.ID,
		UserID:          v.UserID,
		PaymentMethodID: pmID,
		Amount:          v.Amount,
		Currency:        v.Currency,
		Status:          entities.PaymentStatus(v.Status),
		TransactionID:   v.TransactionID,
		ExtraMetadata:   fromPgJSON(v.ExtraMetadata),
		CreatedAt:       v.CreatedAt.Time,
	}
}
