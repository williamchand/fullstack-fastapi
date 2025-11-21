package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type oauthRepository struct {
	queries *dbgen.Queries
}

func NewOAuthRepository(queries *dbgen.Queries) repositories.OAuthRepository {
	return &oauthRepository{queries: queries}
}

func (r *oauthRepository) CreateOAuthAccount(ctx context.Context, oauth *entities.OAuthAccount) error {
	params := dbgen.CreateOAuthAccountParams{
		UserID:         oauth.UserID,
		Provider:       oauth.Provider,
		ProviderUserID: oauth.ProviderUserID,
		AccessToken:    toPgText(oauth.AccessToken),
		RefreshToken:   toPgText(oauth.RefreshToken),
		TokenExpiresAt: toPgTimestamptz(oauth.TokenExpiresAt),
	}

	dbOAuth, err := r.queries.CreateOAuthAccount(ctx, params)
	if err != nil {
		return err
	}

	// Update the entity with generated values
	oauth.ID = dbOAuth.ID
	oauth.CreatedAt = dbOAuth.CreatedAt.Time
	oauth.UpdatedAt = dbOAuth.UpdatedAt.Time

	return nil
}

func (r *oauthRepository) GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*entities.OAuthAccount, error) {
	dbOAuth, err := r.queries.GetOAuthAccount(ctx, dbgen.GetOAuthAccountParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if err != nil {
		return nil, err
	}

	return r.toEntity(&dbOAuth), nil
}

// Additional methods that might be useful

func (r *oauthRepository) GetOAuthAccountByID(ctx context.Context, id uuid.UUID) (*entities.OAuthAccount, error) {
	dbOAuth, err := r.queries.GetOAuthAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.toEntity(&dbOAuth), nil
}

func (r *oauthRepository) GetOAuthAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.OAuthAccount, error) {
	dbOAuths, err := r.queries.GetOAuthAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	oauthAccounts := make([]*entities.OAuthAccount, len(dbOAuths))
	for i, dbOAuth := range dbOAuths {
		oauthAccounts[i] = r.toEntity(&dbOAuth)
	}

	return oauthAccounts, nil
}

func (r *oauthRepository) UpdateOAuthAccountTokens(ctx context.Context, id, userID uuid.UUID, accessToken, refreshToken *string, tokenExpiresAt *time.Time) error {
	params := dbgen.UpdateOAuthAccountTokensParams{
		ID:             id,
		UserID:         userID,
		AccessToken:    toPgText(accessToken),
		RefreshToken:   toPgText(refreshToken),
		TokenExpiresAt: toPgTimestamptz(tokenExpiresAt),
	}

	_, err := r.queries.UpdateOAuthAccountTokens(ctx, params)
	return err
}

func (r *oauthRepository) DeleteOAuthAccount(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.DeleteOAuthAccount(ctx, dbgen.DeleteOAuthAccountParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *oauthRepository) toEntity(dbOAuth *dbgen.OauthAccount) *entities.OAuthAccount {
	return &entities.OAuthAccount{
		ID:             dbOAuth.ID,
		UserID:         dbOAuth.UserID,
		Provider:       dbOAuth.Provider,
		ProviderUserID: dbOAuth.ProviderUserID,
		AccessToken:    fromPgText(dbOAuth.AccessToken),
		RefreshToken:   fromPgText(dbOAuth.RefreshToken),
		TokenExpiresAt: fromPgTime(dbOAuth.TokenExpiresAt),
		CreatedAt:      dbOAuth.CreatedAt.Time,
		UpdatedAt:      dbOAuth.UpdatedAt.Time,
	}
}
