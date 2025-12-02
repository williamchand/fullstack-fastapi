package repositories

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
    TxProvider[UserRepository]

    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    GetByPhone(ctx context.Context, phone string) (*entities.User, error)
    Create(ctx context.Context, user *entities.User) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) (*entities.User, error)
    SetUserRoles(ctx context.Context, userID uuid.UUID, roles []entities.RoleEnum) error
    SetPhoneVerified(ctx context.Context, userID uuid.UUID) error
    SetEmailVerified(ctx context.Context, userID uuid.UUID) error
}
