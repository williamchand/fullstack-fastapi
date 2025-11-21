package database

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type userRepository struct {
    WithTx(tx pgx.Tx) UserRepository
	queries *dbgen.Queries
}

func NewUserRepository(queries *dbgen.Queries) repositories.UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.toEntity(&dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return r.toEntity(&dbUser), nil
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	params := dbgen.CreateUserParams{
		Email:          user.Email,
		PhoneNumber:    toPgText(user.PhoneNumber),
		FullName:       toPgText(user.FullName),
		HashedPassword: toPgText(user.HashedPassword),
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time

	return nil
}

func (r *userRepository) toEntity(dbUser *dbgen.User) *entities.User {
	return &entities.User{
		ID:              dbUser.ID,
		Email:           dbUser.Email,
		PhoneNumber:     fromPgText(dbUser.PhoneNumber),
		FullName:        fromPgText(dbUser.FullName),
		HashedPassword:  fromPgText(dbUser.HashedPassword),
		IsActive:        dbUser.IsActive,
		IsSuperuser:     dbUser.IsSuperuser,
		IsEmailVerified: dbUser.IsEmailVerified,
		IsPhoneVerified: dbUser.IsPhoneVerified,
		IsTOTPEnabled:   dbUser.IsTotpEnabled,
		TOTPSecret:      fromPgText(dbUser.TotpSecret),
		CreatedAt:       dbUser.CreatedAt.Time,
		UpdatedAt:       dbUser.UpdatedAt.Time,
		LastLoginAt:     fromPgTime(dbUser.LastLoginAt),
	}
}
