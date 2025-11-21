package app

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type Repositories struct {
	UserRepo  repositories.UserRepository
	OAuthRepo repositories.OAuthRepository
}

func initRepositories(ctx context.Context, dbURL string) (*Repositories, repositories.ConnectionPool, error) {
	dbPool, err := database.NewPool(ctx, dbURL)
	queries := dbgen.New(dbPool)

	return &Repositories{
		UserRepo:  database.NewUserRepository(queries),
		OAuthRepo: database.NewOAuthRepository(queries),
	}, dbPool, err
}
