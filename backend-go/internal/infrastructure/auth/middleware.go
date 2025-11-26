package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Exact public paths that are always allowed
var (
	publicExactPaths = map[string]map[string]bool{}

	// Public URL prefixes allowed without auth
	publicPrefixes  = []string{}
	publicGRPCExact = map[string]bool{
		"/salonapp.v1.UserService/CreateUser": true,
		"/salonapp.v1.UserService/LoginUser":  true,
	}
	publicGRPCPrefixes = []string{
		"/salonapp.v1.PublicService/",
		"/salonapp.v1.OAuthService/",
	}
	grpcRoleRules = map[string][]string{
		"/salonapp.v1.UserService/GetUser": {string(entities.RoleSuperuser)},
	}
)

type AuthMiddleware struct {
	jwtRepository repositories.JWTRepository
	roleValidator *RoleValidator
	userRepo      repositories.UserRepository // interface to get user by ID
}

func NewAuthMiddleware(jwtRepository repositories.JWTRepository, roleValidator *RoleValidator, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtRepository: jwtRepository,
		roleValidator: roleValidator,
		userRepo:      userRepo,
	}
}

// HTTP middleware
func (m *AuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicHTTPPath(r.URL.Path, r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		token := extractTokenFromHeader(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtRepository.ValidateToken(token)
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

	claims, err := m.jwtRepository.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	user, err := m.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || !user.IsActive {
		return nil, status.Error(codes.Unauthenticated, "user not found or inactive")
	}

	// ROLE AUTHORIZATION (NEW)
	requiredRoles := m.roleValidator.RequiredGRPCRoles(info.FullMethod)
	if !m.roleValidator.HasRole(user, requiredRoles...) {
		return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
	}

	// Add user to context
	ctx = WithUser(ctx, user)
	return handler(ctx, req)
}

func (r *RoleValidator) RequiredGRPCRoles(method string) []string {
	return grpcRoleRules[method]
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// gRPC metadata keys are lowercase
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return ""
	}

	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeaders[0], " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func isPublicMethod(method string) bool {
	if _, ok := publicGRPCExact[method]; ok {
		return true
	}

	// 2. Prefix match
	for _, prefix := range publicGRPCPrefixes {
		if strings.HasPrefix(method, prefix) {
			return true
		}
	}

	return false
}

func isPublicHTTPPath(path string, method string) bool {
	if publicPath, ok := publicExactPaths[path]; ok {
		if _, ok := publicPath[method]; ok {
			return true
		}
		return false
	}

	for _, prefix := range publicPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
