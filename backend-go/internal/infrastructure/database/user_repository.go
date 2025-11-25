package database

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	queries *dbgen.Queries
	db      repositories.ConnectionPool
}

func NewUserRepository(queries *dbgen.Queries, db repositories.ConnectionPool) repositories.UserRepository {
	return &userRepository{queries: queries, db: db}
}

// WithTx sets the transaction for the repository
func (r *userRepository) WithTx(tx pgx.Tx) repositories.UserRepository {
	return &userRepository{
		queries: r.queries.WithTx(tx),
		db:      r.db,
	}
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
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return r.toEntity(&dbUser), nil
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	params := dbgen.CreateUserParams{
		Email:          user.Email,
		PhoneNumber:    toPgText(user.PhoneNumber),
		FullName:       toPgText(user.FullName),
		HashedPassword: toPgText(user.HashedPassword),
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return r.toEntity(&dbUser), err
	}

	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time

	return r.toEntity(&dbUser), err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	params := dbgen.CreateUserParams{
		Email:          user.Email,
		PhoneNumber:    toPgText(user.PhoneNumber),
		FullName:       toPgText(user.FullName),
		HashedPassword: toPgText(user.HashedPassword),
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return r.toEntity(&dbUser), err
	}

	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time

	// Implementation depends on your update strategy
	// This could execute the update or return a transactional version
	return r.toEntity(&dbUser), err
}

func (r *userRepository) GetRoles(ctx context.Context, roles []entities.RoleEnum) ([]int32, error) {
	roleNames := []string{}
	for _, r := range roles {
		roleNames = append(roleNames, string(r))
	}
	dbRoles, err := r.queries.GetRole(ctx, roleNames)
	if err != nil {
		return nil, err
	}

	return dbRoles, nil
}

func (r *userRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entities.Role, error) {
	dbRoles, err := r.queries.GetUserRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles := make([]entities.Role, len(dbRoles))
	for i, dbRole := range dbRoles {
		roles[i] = entities.Role{
			ID:          dbRole.ID,
			Name:        dbRole.Name,
			Description: fromPgText(dbRole.Description),
			CreatedAt:   dbRole.CreatedAt.Time,
			UpdatedAt:   dbRole.UpdatedAt.Time,
		}
	}

	return roles, nil
}

func (r *userRepository) AssignRole(ctx context.Context, userID uuid.UUID, roleID int32) error {
	params := dbgen.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	}
	return r.queries.AssignRoleToUser(ctx, params)
}

func (r *userRepository) SetUserRoles(ctx context.Context, userID uuid.UUID, roles []entities.RoleEnum) error {
	roleIDs, err := r.GetRoles(ctx, roles)
	if err != nil {
		return err
	}

	err = r.queries.DeleteUserRole(ctx, userID)
	if err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		err = r.AssignRole(ctx, userID, roleID)
		if err != nil {
			return err
		}
	}

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
		IsEmailVerified: dbUser.IsEmailVerified,
		IsPhoneVerified: dbUser.IsPhoneVerified,
		IsTOTPEnabled:   dbUser.IsTotpEnabled,
		TOTPSecret:      fromPgText(dbUser.TotpSecret),
		CreatedAt:       dbUser.CreatedAt.Time,
		UpdatedAt:       dbUser.UpdatedAt.Time,
		LastLoginAt:     fromPgTime(dbUser.LastLoginAt),
	}
}
