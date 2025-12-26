package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type VerificationCodeRepository interface {
	TxProvider[VerificationCodeRepository]

	Create(ctx context.Context, v *entities.VerificationCode) error
	CreateNoUser(ctx context.Context, code string, vType entities.VerificationType, extraMetadata map[string]any, expiresAt time.Time) (*entities.VerificationCode, error)
	GetLatestUnused(ctx context.Context, userID uuid.UUID, vType entities.VerificationType) (*entities.VerificationCode, error)
	GetByCode(ctx context.Context, userID uuid.UUID, vType entities.VerificationType, code string) (*entities.VerificationCode, error)
	GetByCodeOnly(ctx context.Context, vType entities.VerificationType, code string) (*entities.VerificationCode, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
}
