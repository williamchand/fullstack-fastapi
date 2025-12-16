package util

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

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
