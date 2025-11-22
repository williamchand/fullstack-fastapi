package services

import (
	"context"
	"errors"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidRole     = errors.New("invalid role")
	ErrNoRolesProvided = errors.New("no roles provided")
)

type UserService struct {
	userRepo  repositories.UserRepository
	oauthRepo repositories.OAuthRepository
}

func NewUserService(userRepo repositories.UserRepository, oauthRepo repositories.OAuthRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
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

// Role Management Methods

// SetUserRoles sets/replaces all roles for a user
func (s *UserService) SetUserRoles(ctx context.Context, userID string, roleIDs []int32) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.SetUserRoles(ctx, userUUID, roleIDs)
}

// GetUserRoles retrieves all roles for a user
func (s *UserService) GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.userRepo.GetUserRoles(ctx, userUUID)
}

// AssignRoleToUser assigns a single role to a user
func (s *UserService) AssignRoleToUser(ctx context.Context, userID string, roleID int32) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.AssignRole(ctx, userUUID, roleID)
}

// RemoveRoleFromUser removes a specific role from a user
func (s *UserService) RemoveRoleFromUser(ctx context.Context, userID string, roleID int32) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// Get current roles
	currentRoles, err := s.userRepo.GetUserRoles(ctx, userUUID)
	if err != nil {
		return err
	}

	// Filter out the role to remove
	var newRoleIDs []int32
	for _, role := range currentRoles {
		if role.ID != roleID {
			newRoleIDs = append(newRoleIDs, role.ID)
		}
	}

	// Set the updated roles
	return s.userRepo.SetUserRoles(ctx, userUUID, newRoleIDs)
}

// HasRole checks if a user has a specific role
func (s *UserService) HasRole(ctx context.Context, userID string, roleID int32) (bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, ErrUserNotFound
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userUUID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.ID == roleID {
			return true, nil
		}
	}

	return false, nil
}

// HasAnyRole checks if a user has any of the specified roles
func (s *UserService) HasAnyRole(ctx context.Context, userID string, roleIDs []int32) (bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, ErrUserNotFound
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userUUID)
	if err != nil {
		return false, err
	}

	roleIDSet := make(map[int32]bool)
	for _, roleID := range roleIDs {
		roleIDSet[roleID] = true
	}

	for _, role := range roles {
		if roleIDSet[role.ID] {
			return true, nil
		}
	}

	return false, nil
}

// ClearUserRoles removes all roles from a user
func (s *UserService) ClearUserRoles(ctx context.Context, userID string) error {
	return s.SetUserRoles(ctx, userID, []int32{})
}
