package services

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  repositories.UserRepository
	oauthRepo repositories.OAuthRepository
	jwtRepo   repositories.JWTRepository
}

func NewUserService(userRepo repositories.UserRepository, oauthRepo repositories.OAuthRepository, jwtRepo repositories.JWTRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
		jwtRepo:   jwtRepo,
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

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, email, password, fullName, phoneNumber string) (*entities.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	// Hash password
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

	user, err = s.userRepo.Update(ctx, user)
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
