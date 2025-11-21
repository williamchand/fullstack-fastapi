package repositories

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) UserRepository
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entities.Role, error)
	AssignRole(ctx context.Context, userID uuid.UUID, roleID int32) error
}

type OAuthRepository interface {
	CreateOAuthAccount(ctx context.Context, oauth *entities.OAuthAccount) error
	GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*entities.OAuthAccount, error)
}
