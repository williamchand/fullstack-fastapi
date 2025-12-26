package app

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	dokunfra "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/doku"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/smtp"
	stripeinfra "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/stripe"
	wahainfra "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/waha"
)

type AppServices struct {
	UserService    *services.UserService
	OauthService   *services.OAuthService
	BillingService *services.BillingService
}

func initServices(cfg *config.Config, repo *Repositories) (*AppServices, error) {
	jwtService, err := jwt.NewService(cfg)
	if err != nil {
		return nil, err
	}
	smtpSender := smtp.NewSMTPSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From)
	wahaClient := wahainfra.New(cfg.WAHA.URL, cfg.WAHA.APIKey, cfg.WAHA.Session)
	stripeClient := stripeinfra.New(cfg.Stripe.SecretKey)
	dokuClient := dokunfra.New(cfg.Doku.BaseURL, cfg.Doku.ClientID, cfg.Doku.SecretKey)
	return &AppServices{
		UserService:    services.NewUserService(cfg, repo.UserRepo, repo.OAuthRepo, repo.TransactionManager, jwtService, repo.EmailTemplateRepo, repo.VerificationRepo, smtpSender, wahaClient),
		OauthService:   services.NewOAuthService(cfg.GetOauthConfig(), repo.OAuthRepo, repo.UserRepo, repo.TransactionManager, jwtService),
		BillingService: services.NewBillingService(cfg, repo.SubscriptionRepo, repo.PaymentRepo, stripeClient, dokuClient),
	}, nil
}
