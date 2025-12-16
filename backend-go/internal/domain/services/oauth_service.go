package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ProviderUserInfo represents the user information returned by Google OAuth
type ProviderUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type OAuthConfigService struct {
	configService *oauth2.Config
	infoURL       string
}

type OAuthService struct {
	config    map[string]OAuthConfigService
	oauthRepo repositories.OAuthRepository
	userRepo  repositories.UserRepository
	txManager repositories.TransactionManager
	jwtRepo   repositories.JWTRepository
}

func NewOAuthService(
	cfg *config.OAuthConfig,
	oauthRepo repositories.OAuthRepository,
	userRepo repositories.UserRepository,
	txManager repositories.TransactionManager,
	jwtRepo repositories.JWTRepository,
) *OAuthService {
	config := map[string]OAuthConfigService{}
	config["google"] = OAuthConfigService{
		configService: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		infoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
	return &OAuthService{
		config:    config,
		oauthRepo: oauthRepo,
		userRepo:  userRepo,
		txManager: txManager,
		jwtRepo:   jwtRepo,
	}
}

func (s *OAuthService) GetAuthURL(provider string) (string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", "", err
	}

	return s.config[provider].configService.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent")), state, nil
}

func (s *OAuthService) HandleCallback(ctx context.Context, provider string, code string) (*entities.TokenPair, error) {
	token, err := s.config[provider].configService.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	userInfo, err := s.getUserInfo(ctx, provider, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	var user *entities.User
	var isNewUser bool
	// Execute in transaction to ensure data consistency
	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		// Use repository with transaction
		oauthRepoTx := s.oauthRepo.WithTx(tx)
		userRepoTx := s.userRepo.WithTx(tx)

		// Check if OAuth account exists
		oauthAccount, err := oauthRepoTx.GetOAuthAccount(ctx, provider, userInfo.ID)
		if err == nil && oauthAccount != nil {
			// Update OAuth tokens
			err = oauthRepoTx.UpdateOAuthAccountTokens(ctx, oauthAccount.ID, oauthAccount.UserID, &token.AccessToken, &token.RefreshToken, &token.Expiry)
			if err != nil {
				return fmt.Errorf("failed to update OAuth tokens: %w", err)
			}

			// Get the user
			user, err = userRepoTx.GetByID(ctx, oauthAccount.UserID)
			return err
		}

		// New user - create account
		user, err = s.createUserFromOAuth(ctx, userRepoTx, oauthRepoTx, provider, userInfo, token)
		isNewUser = true
		return err
	})

	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtRepo.GenerateToken(user.ID, user.Email, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtRepo.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return &entities.TokenPair{
		User:             user,
		AccessToken:      accessToken.Token,
		RefreshToken:     refreshToken.Token,
		ExpiresAt:        accessToken.ExpiresAt,
		RefreshExpiresAt: refreshToken.ExpiresAt,
		IsNewUser:        isNewUser,
	}, nil
}

// getUserInfo retrieves user information from Google API
func (s *OAuthService) getUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*ProviderUserInfo, error) {
	client := s.config[provider].configService.Client(ctx, token)

	resp, err := client.Get(s.config[provider].infoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo ProviderUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// createUserFromOAuth creates a new user from OAuth information
func (s *OAuthService) createUserFromOAuth(
	ctx context.Context,
	userRepo repositories.UserRepository,
	oauthRepo repositories.OAuthRepository,
	provider string,
	userInfo *ProviderUserInfo,
	token *oauth2.Token,
) (*entities.User, error) {
	// Check if a user with this email already exists
	existingUser, err := userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil && existingUser != nil {
		// User exists with this email, link OAuth account to existing user
		oauthAccount := &entities.OAuthAccount{
			Provider:       provider,
			ProviderUserID: userInfo.ID,
			UserID:         existingUser.ID,
			AccessToken:    &token.AccessToken,
			RefreshToken:   &token.RefreshToken,
			TokenExpiresAt: &token.Expiry,
			ProviderData:   s.buildProviderData(userInfo),
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
		Roles:           []string{string(entities.RoleCustomer)},
	}

	// Create the user in database
	user, err = userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Set user roles
	err = userRepo.SetUserRoles(ctx, user.ID, []entities.RoleEnum{entities.RoleCustomer})
	if err != nil {
		return nil, fmt.Errorf("failed to set user roles: %w", err)
	}

	// Create OAuth account
	oauthAccount := &entities.OAuthAccount{
		Provider:       provider,
		ProviderUserID: userInfo.ID,
		UserID:         user.ID,
		AccessToken:    &token.AccessToken,
		RefreshToken:   &token.RefreshToken,
		TokenExpiresAt: &token.Expiry,
		ProviderData:   s.buildProviderData(userInfo),
	}

	err = oauthRepo.CreateOAuthAccount(ctx, oauthAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth account: %w", err)
	}

	return user, nil
}

// buildProviderData builds the provider data JSON from user info
func (s *OAuthService) buildProviderData(userInfo *ProviderUserInfo) map[string]any {
	return map[string]any{
		"name":        userInfo.Name,
		"given_name":  userInfo.GivenName,
		"family_name": userInfo.FamilyName,
		"picture":     userInfo.Picture,
		"locale":      userInfo.Locale,
	}
}

// RefreshToken refreshes the OAuth token if it's expired
func (s *OAuthService) RefreshToken(ctx context.Context, provider string, userID string) error {
	return s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		oauthRepo := s.oauthRepo.WithTx(tx)

		// Get the OAuth account for this user
		oauthAccounts, err := oauthRepo.GetOAuthAccount(ctx, provider, userID)
		if err != nil {
			return fmt.Errorf("failed to get OAuth accounts: %w", err)
		}

		if oauthAccounts == nil {
			return fmt.Errorf("no Google OAuth account found for user")
		}

		// Create token source with refresh token
		token := &oauth2.Token{
			AccessToken:  *oauthAccounts.AccessToken,
			RefreshToken: *oauthAccounts.RefreshToken,
			Expiry:       *oauthAccounts.TokenExpiresAt,
		}

		newToken, err := s.config[provider].configService.TokenSource(ctx, token).Token()
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update the tokens in database
		return oauthRepo.UpdateOAuthAccountTokens(ctx, oauthAccounts.ID, oauthAccounts.UserID, &newToken.AccessToken, &newToken.RefreshToken, &newToken.Expiry)
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
