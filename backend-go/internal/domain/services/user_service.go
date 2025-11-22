package services

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   repositories.UserRepository
	oauthRepo  repositories.OAuthRepository
	jwtService jwt.JWTService
}

func NewUserService(userRepo repositories.UserRepository, oauthRepo repositories.OAuthRepository, jwtService jwt.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		oauthRepo:  oauthRepo,
		jwtService: jwtService,
	}
}

// Existing methods...

func (s *UserService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) CreateUser(ctx context.Context, email, password, fullName, phoneNumber string) (*entities.User, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	user := &entities.User{
		Email:          email,
		HashedPassword: &hashedPasswordStr,
		FullName:       &fullName,
		IsActive:       true,
	}

	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, email, password, fullName, phoneNumber string) (*entities.User, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	user := &entities.User{
		Email:          email,
		HashedPassword: &hashedPasswordStr,
		FullName:       &fullName,
		IsActive:       true,
	}

	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ValidatePassword(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	if !user.IsActive {
		return nil, ErrInvalidPassword
	}

	if user.HashedPassword == nil {
		return nil, ErrInvalidPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
