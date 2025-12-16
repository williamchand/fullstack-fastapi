package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
)

type Middleware struct {
	Auth *auth.AuthMiddleware
}

func initMiddleware(cfg *config.Config, repo *Repositories) (*Middleware, error) {
	jwtService, err := jwt.NewService(cfg)
	if err != nil {
		return nil, err
	}
	return &Middleware{
		Auth: auth.NewAuthMiddleware(jwtService, repo.UserRepo),
	}, nil
}
