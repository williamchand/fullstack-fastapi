package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
)

type AppServices struct {
	UserService  *services.UserService
	OauthService *services.OAuthService
}

func initServices(cfg *config.Config, repo *Repositories) *AppServices {
	return &AppServices{
		UserService:  services.NewUserService(repo.UserRepo, repo.OAuthRepo),
		OauthService: services.NewOAuthService(cfg.OAuth, repo.OAuthRepo, repo.UserRepo, repo.TransactionManager),
	}
}
