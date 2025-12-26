package app

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type Repositories struct {
	TransactionManager repositories.TransactionManager
	UserRepo           repositories.UserRepository
	OAuthRepo          repositories.OAuthRepository
	EmailTemplateRepo  repositories.EmailTemplateRepository
	VerificationRepo   repositories.VerificationCodeRepository
	SubscriptionRepo   repositories.SubscriptionRepository
	PaymentRepo        repositories.PaymentRepository
}

func initRepositories(ctx context.Context, dbURL string) (*Repositories, repositories.ConnectionPool, error) {
	dbPool, err := database.NewPool(ctx, dbURL)
	if err != nil {
		return nil, nil, err
	}
	queries := dbgen.New(dbPool)

	return &Repositories{
		TransactionManager: database.NewTransactionManager(dbPool),
		UserRepo:           database.NewUserRepository(queries, dbPool),
		OAuthRepo:          database.NewOAuthRepository(queries, dbPool),
		EmailTemplateRepo:  database.NewEmailTemplateRepository(queries, dbPool),
		VerificationRepo:   database.NewVerificationCodeRepository(queries, dbPool),
		SubscriptionRepo:   database.NewSubscriptionRepository(queries, dbPool),
		PaymentRepo:        database.NewPaymentRepository(queries, dbPool),
	}, dbPool, err
}
