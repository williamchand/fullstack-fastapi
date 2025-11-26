package main

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
)

var (
	cfg         *config.Config
	userService *services.UserService
)

func main() {
	setup()
	generateTestAccounts()
}

func setup() {
	cfg, _ = config.Load()
	ctx := context.Background()
	dbPool, _ := database.NewPool(ctx, cfg.DatabaseURL)
	queries := dbgen.New(dbPool)
	transactionManager := database.NewTransactionManager(dbPool)
	userRepo := database.NewUserRepository(queries, dbPool)
	oAuthRepo := database.NewOAuthRepository(queries, dbPool)
	jwtService, _ := jwt.NewService(cfg)
	userService = services.NewUserService(userRepo, oAuthRepo, transactionManager, jwtService)
}

func generateTestAccounts() {
	ctx := context.Background()
	fullName := "Superuser"
	userService.CreateUser(ctx, cfg.Superuser.Username, cfg.Superuser.Password, fullName, "", []entities.RoleEnum{entities.RoleSuperuser}, true)
}
