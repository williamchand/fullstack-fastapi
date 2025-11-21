package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GoogleOAuthService struct {
	config    *oauth2.Config
	oauthRepo repositories.OAuthRepository
	userRepo  repositories.UserRepository
}

func NewGoogleOAuthService(cfg *GoogleOAuthConfig, oauthRepo repositories.OAuthRepository, userRepo repositories.UserRepository) *GoogleOAuthService {
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthService{
		config:    config,
		oauthRepo: oauthRepo,
		userRepo:  userRepo,
	}
}

func (s *GoogleOAuthService) GetAuthURL() (string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", "", err
	}

	return s.config.AuthCodeURL(state), state, nil
}

func (s *GoogleOAuthService) HandleCallback(ctx context.Context, code string) (*entities.User, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Get user info from Google
	client := s.config.Client(ctx, token)
	userInfo, err := s.getUserInfo(client)
	if err != nil {
		return nil, err
	}

	// Check if OAuth account exists
	oauthAccount, err := s.oauthRepo.GetOAuthAccount(ctx, "google", userInfo.ID)
	if err == nil && oauthAccount != nil {
		// Existing user - update tokens and return user
		return s.userRepo.GetByID(ctx, oauthAccount.UserID)
	}

	// New user - create account
	return s.createUserFromOAuth(ctx, userInfo, token)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
