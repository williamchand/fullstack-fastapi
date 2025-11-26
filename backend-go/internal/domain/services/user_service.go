package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  repositories.UserRepository
	oauthRepo repositories.OAuthRepository
	txManager repositories.TransactionManager
	jwtRepo   repositories.JWTRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	oauthRepo repositories.OAuthRepository,
	txManager repositories.TransactionManager,
	jwtRepo repositories.JWTRepository,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
		txManager: txManager,
		jwtRepo:   jwtRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) CreateUser(ctx context.Context, email, password, fullName, phoneNumber string, roles []entities.RoleEnum, isEmailVerified bool) (*entities.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	user := &entities.User{
		Email:           email,
		HashedPassword:  &hashedPasswordStr,
		FullName:        &fullName,
		IsActive:        true,
		IsEmailVerified: isEmailVerified,
	}

	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)

		user, err = userRepoTx.Create(ctx, user)
		if err != nil {
			return err
		}

		err = userRepoTx.SetUserRoles(ctx, user.ID, roles)
		if err != nil {
			return err
		}
		return nil
	})

	for _, role := range roles {
		user.Roles = append(user.Roles, string(role))
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, email, password, fullName, phoneNumber string, roles []entities.RoleEnum) (*entities.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser == nil {
		return nil, ErrUserNotFound
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:          existingUser.ID,
		Email:       email,
		FullName:    &fullName,
		PhoneNumber: &phoneNumber,
		IsActive:    true,
	}
	hashedPasswordStr := string(hashedPassword)
	if password != "" {
		user.HashedPassword = &hashedPasswordStr
	}

	err = s.txManager.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		userRepoTx := s.userRepo.WithTx(tx)

		user, err = userRepoTx.Update(ctx, user)
		if err != nil {
			return err
		}

		err = userRepoTx.SetUserRoles(ctx, user.ID, roles)
		if err != nil {
			return err
		}
		return nil
	})

	for _, role := range roles {
		user.Roles = append(user.Roles, string(role))
	}

	return user, nil
}

func (s *UserService) ValidatePassword(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}
	if !user.IsEmailVerified {
		return nil, ErrInvalidEmailNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
func (s *UserService) Login(
	ctx context.Context,
	username string,
	password string,
) (*entities.TokenPair, error) {
	user, err := s.ValidatePassword(ctx, username, password)
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
		User:         user,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt,
		IsNewUser:    false,
	}, nil
}
