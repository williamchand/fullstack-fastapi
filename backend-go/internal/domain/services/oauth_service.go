package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
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

// GoogleUserInfo represents the user information returned by Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleOAuthService struct {
	config    *oauth2.Config
	oauthRepo repositories.OAuthRepository
	userRepo  repositories.UserRepository
	txManager repositories.TransactionManager
}

func NewGoogleOAuthService(
	cfg *GoogleOAuthConfig,
	oauthRepo repositories.OAuthRepository,
	userRepo repositories.UserRepository,
	txManager repositories.TransactionManager,
) *GoogleOAuthService {
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
		txManager: txManager,
	}
}

func (s *GoogleOAuthService) GetAuthURL() (string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", "", err
	}

	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline), state, nil
}

func (s *GoogleOAuthService) HandleCallback(ctx context.Context, code string) (*entities.User, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	userInfo, err := s.getUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	var user *entities.User

	// Execute in transaction to ensure data consistency
	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		// Use repository with transaction
		oauthRepoTx := s.oauthRepo.WithTx(tx)
		userRepoTx := s.userRepo.WithTx(tx)

		// Check if OAuth account exists
		oauthAccount, err := oauthRepoTx.GetOAuthAccount(ctx, "google", userInfo.ID)
		if err == nil && oauthAccount != nil {
			// Update OAuth tokens
			err = oauthRepoTx.UpdateOAuthTokens(ctx, oauthAccount.ID, token.AccessToken, token.RefreshToken, token.Expiry)
			if err != nil {
				return fmt.Errorf("failed to update OAuth tokens: %w", err)
			}

			// Get the user
			user, err = userRepoTx.GetByID(ctx, oauthAccount.UserID)
			return err
		}

		// New user - create account
		user, err = s.createUserFromOAuth(ctx, userRepoTx, oauthRepoTx, userInfo, token)
		return err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// getUserInfo retrieves user information from Google API
func (s *GoogleOAuthService) getUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// createUserFromOAuth creates a new user from OAuth information
func (s *GoogleOAuthService) createUserFromOAuth(
	ctx context.Context,
	userRepo repositories.UserRepository,
	oauthRepo repositories.OAuthRepository,
	userInfo *GoogleUserInfo,
	token *oauth2.Token,
) (*entities.User, error) {
	// Check if a user with this email already exists
	existingUser, err := userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil && existingUser != nil {
		// User exists with this email, link OAuth account to existing user
		oauthAccount := &entities.OAuthAccount{
			Provider:     "google",
			ProviderID:   userInfo.ID,
			UserID:       existingUser.ID,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenExpiry:  token.Expiry,
			ProviderData: s.buildProviderData(userInfo),
		}

		err = oauthRepo.CreateOAuthAccount(ctx, oauthAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to create OAuth account for existing user: %w", err)
		}

		return existingUser, nil
	}

	// Create new user
	user := &entities.User{
		Email:           userInfo.Email,
		FullName:        &userInfo.Name,
		IsActive:        true,
		IsEmailVerified: userInfo.VerifiedEmail,
	}

	// Create the user in database
	err = userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create OAuth account
	oauthAccount := &entities.OAuthAccount{
		Provider:     "google",
		ProviderID:   userInfo.ID,
		UserID:       user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
		ProviderData: s.buildProviderData(userInfo),
	}

	err = oauthRepo.CreateOAuthAccount(ctx, oauthAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth account: %w", err)
	}

	return user, nil
}

// buildProviderData builds the provider data JSON from user info
func (s *GoogleOAuthService) buildProviderData(userInfo *GoogleUserInfo) map[string]interface{} {
	return map[string]interface{}{
		"name":        userInfo.Name,
		"given_name":  userInfo.GivenName,
		"family_name": userInfo.FamilyName,
		"picture":     userInfo.Picture,
		"locale":      userInfo.Locale,
	}
}

// RefreshToken refreshes the OAuth token if it's expired
func (s *GoogleOAuthService) RefreshToken(ctx context.Context, userID string) error {
	return s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		oauthRepo := s.oauthRepo.WithTx(tx)

		// Get the OAuth account for this user
		oauthAccounts, err := oauthRepo.GetOAuthAccountsByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get OAuth accounts: %w", err)
		}

		var googleAccount *entities.OAuthAccount
		for _, account := range oauthAccounts {
			if account.Provider == "google" {
				googleAccount = &account
				break
			}
		}

		if googleAccount == nil {
			return fmt.Errorf("no Google OAuth account found for user")
		}

		// Create token source with refresh token
		token := &oauth2.Token{
			AccessToken:  googleAccount.AccessToken,
			RefreshToken: googleAccount.RefreshToken,
			Expiry:       googleAccount.TokenExpiry,
		}

		newToken, err := s.config.TokenSource(ctx, token).Token()
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update the tokens in database
		return oauthRepo.UpdateOAuthTokens(ctx, googleAccount.ID, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
	})
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
