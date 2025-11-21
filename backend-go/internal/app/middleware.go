package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"
)

type Middleware struct {
	Auth *auth.AuthMiddleware
}

func initMiddleware(cfg *config.Config, repo *Repositories) *Middleware {
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	roleValidator := auth.NewRoleValidator()

	return &Middleware{
		Auth: auth.NewAuthMiddleware(jwtService, roleValidator, repo.UserRepo),
	}
}
