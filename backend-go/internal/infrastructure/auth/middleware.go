package auth

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthMiddleware struct {
	jwtService    *JWTService
	roleValidator *RoleValidator
	userRepo      UserRepository // interface to get user by ID
}

func NewAuthMiddleware(jwtService *JWTService, roleValidator *RoleValidator, userRepo UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:    jwtService,
		roleValidator: roleValidator,
		userRepo:      userRepo,
	}
}

// HTTP middleware
func (m *AuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractTokenFromHeader(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil || !user.IsActive {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GRPC interceptor for authentication
func (m *AuthMiddleware) GRPCAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip auth for certain methods (like login)
	if isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	// Extract token from context
	token := extractTokenFromGRPCContext(ctx)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	claims, err := m.jwtService.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	user, err := m.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || !user.IsActive {
		return nil, status.Error(codes.Unauthenticated, "user not found or inactive")
	}

	// Add user to context
	ctx = WithUser(ctx, user)
	return handler(ctx, req)
}

// GRPC interceptor for role-based authorization
func (m *AuthMiddleware) GRPCRoleInterceptor(requiredRoles ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		user := UserFromContext(ctx)
		if user == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if !m.roleValidator.HasRole(user, requiredRoles...) {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func extractTokenFromGRPCContext(ctx context.Context) string {
	// This depends on how you pass tokens in gRPC
	// Common approach is using metadata
	return ""
}

func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/user.v1.UserService/CreateUser",
		"/auth.v1.AuthService/Login",
	}

	for _, public := range publicMethods {
		if method == public {
			return true
		}
	}
	return false
}
