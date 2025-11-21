package repositories

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type OAuthRepository interface {
	TxProvider[OAuthRepository]

	CreateOAuthAccount(ctx context.Context, oauth *entities.OAuthAccount) error
	GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*entities.OAuthAccount, error)
}
