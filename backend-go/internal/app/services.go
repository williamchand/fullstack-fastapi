package app

import "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"

type AppServices struct {
	UserService *services.UserService
}

func initServices(repo *Repositories) *AppServices {
	return &AppServices{
		UserService: services.NewUserService(repo.UserRepo, repo.OAuthRepo),
	}
}
