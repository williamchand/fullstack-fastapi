package auth

import (
	"context"
	"strings"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type RoleValidator struct {
}

func NewRoleValidator() *RoleValidator {
	return &RoleValidator{}
}

// HasRole checks if user has any of the required roles
func (v *RoleValidator) HasRole(user *entities.User, requiredRoles ...string) bool {
	if user.IsSuperuser {
		return true
	}

	userRoles := make(map[string]bool)
	for _, role := range user.Roles {
		userRoles[strings.ToLower(role.Name)] = true
	}

	for _, required := range requiredRoles {
		if userRoles[strings.ToLower(required)] {
			return true
		}
	}

	return false
}

// HasAllRoles checks if user has all required roles
func (v *RoleValidator) HasAllRoles(user *entities.User, requiredRoles ...string) bool {
	if user.IsSuperuser {
		return true
	}

	userRoles := make(map[string]bool)
	for _, role := range user.Roles {
		userRoles[strings.ToLower(role.Name)] = true
	}

	for _, required := range requiredRoles {
		if !userRoles[strings.ToLower(required)] {
			return false
		}
	}

	return true
}

// Context key type for safety
type contextKey string

const (
	userContextKey contextKey = "user"
)

// WithUser adds user to context
func WithUser(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves user from context
func UserFromContext(ctx context.Context) *entities.User {
	user, _ := ctx.Value(userContextKey).(*entities.User)
	return user
}
