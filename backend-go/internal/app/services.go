package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/smtp"
	stripeinfra "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/stripe"
)

type AppServices struct {
	UserService       *services.UserService
	OauthService      *services.OAuthService
	DataSourceService *services.DataSourceService
	BillingService    *services.BillingService
	WeddingService    *services.WeddingService
}

func initServices(cfg *config.Config, repo *Repositories) (*AppServices, error) {
	jwtService, err := jwt.NewService(cfg)
	if err != nil {
		return nil, err
	}
	smtpSender := smtp.NewSMTPSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From)
	stripeClient := stripeinfra.New(cfg.Stripe.SecretKey)
	return &AppServices{
		UserService:       services.NewUserService(repo.UserRepo, repo.OAuthRepo, repo.TransactionManager, jwtService, repo.EmailTemplateRepo, repo.VerificationRepo, smtpSender),
		OauthService:      services.NewOAuthService(cfg.GetOauthConfig(), repo.OAuthRepo, repo.UserRepo, repo.TransactionManager, jwtService),
		DataSourceService: services.NewDataSourceService(cfg, repo.DataSourceRepo, repo.AICredRepo, repo.TransactionManager),
		BillingService:    services.NewBillingService(cfg, repo.SubscriptionRepo, repo.PaymentRepo, stripeClient),
		WeddingService:    services.NewWeddingService(repo.WeddingRepo, repo.GuestRepo, repo.TemplateRepo, repo.PaymentRepo, repo.SubscriptionRepo),
	}, nil
}
