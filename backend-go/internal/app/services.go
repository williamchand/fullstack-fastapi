package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
)

type AppServices struct {
	UserService  *services.UserService
	OauthService *services.OAuthService
}

func initServices(cfg *config.Config, repo *Repositories) (*AppServices, error) {
	jwtService, err := jwt.NewService(cfg)
	if err != nil {
		return nil, err
	}
	return &AppServices{
		UserService:  services.NewUserService(repo.UserRepo, repo.OAuthRepo, repo.TransactionManager, jwtService),
		OauthService: services.NewOAuthService(cfg.GetOauthConfig(), repo.OAuthRepo, repo.UserRepo, repo.TransactionManager, jwtService),
	}, nil
}
