package database

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type verificationCodeRepository struct {
    queries *dbgen.Queries
    db      repositories.ConnectionPool
}

func NewVerificationCodeRepository(queries *dbgen.Queries, db repositories.ConnectionPool) repositories.VerificationCodeRepository {
    return &verificationCodeRepository{queries: queries, db: db}
}

func (r *verificationCodeRepository) WithTx(tx pgx.Tx) repositories.VerificationCodeRepository {
    return &verificationCodeRepository{
        queries: r.queries.WithTx(tx),
        db:      r.db,
    }
}

func (r *verificationCodeRepository) Create(ctx context.Context, v *entities.VerificationCode) error {
    params := dbgen.CreateVerificationCodeParams{
        UserID:           v.UserID,
        VerificationCode: v.Code,
        VerificationType: string(v.Type),
        ExtraMetadata:    toPgJSON(v.ExtraMetadata),
        ExpiresAt:        toPgTime(v.ExpiresAt),
    }
    res, err := r.queries.CreateVerificationCode(ctx, params)
    if err != nil {
        return err
    }
    v.ID = res.ID
    v.CreatedAt = res.CreatedAt.Time
    return nil
}

func (r *verificationCodeRepository) GetLatestUnused(ctx context.Context, userID uuid.UUID, vType entities.VerificationType) (*entities.VerificationCode, error) {
    res, err := r.queries.GetLatestUnusedVerificationCode(ctx, dbgen.GetLatestUnusedVerificationCodeParams{
        UserID:           userID,
        VerificationType: string(vType),
    })
    if err != nil {
        return nil, err
    }
    return r.toEntity(&res), nil
}

func (r *verificationCodeRepository) GetByCode(ctx context.Context, userID uuid.UUID, vType entities.VerificationType, code string) (*entities.VerificationCode, error) {
    res, err := r.queries.GetVerificationCodeByCode(ctx, dbgen.GetVerificationCodeByCodeParams{
        UserID:           userID,
        VerificationType: string(vType),
        VerificationCode: code,
    })
    if err != nil {
        return nil, err
    }
    return r.toEntity(&res), nil
}

func (r *verificationCodeRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
    _, err := r.queries.MarkVerificationCodeUsed(ctx, id)
    return err
}

func (r *verificationCodeRepository) toEntity(v *dbgen.VerificationCode) *entities.VerificationCode {
    return &entities.VerificationCode{
        ID:            v.ID,
        UserID:        v.UserID,
        Code:          v.VerificationCode,
        Type:          entities.VerificationType(v.VerificationType),
        ExtraMetadata: fromPgJSON(v.ExtraMetadata),
        CreatedAt:     v.CreatedAt.Time,
        ExpiresAt:     v.ExpiresAt.Time,
        UsedAt:        fromPgTime(v.UsedAt),
    }
}