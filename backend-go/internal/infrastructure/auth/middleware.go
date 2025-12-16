package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Exact public paths that are always allowed
var (
	publicExactPaths = map[string]map[string]bool{
		"/v1/login/access-token":     {"POST": true},
		"/v1/login/refresh-token":    {"POST": true},
		"/v1/user/verify-email":      {"POST": true},
		"/v1/user/resend-email":      {"POST": true},
		"/v1/user":                   {"POST": true},
		"/v1/password-recovery":      {"POST": true},
		"/v1/reset-password":         {"POST": true},
		"/v1/login/phone":            {"POST": true},
		"/v1/user/register-phone":    {"POST": true},
		"/v1/user/request-phone-otp": {"POST": true},
		"/v1/user/verify-phone-otp":  {"POST": true},
	}

	// Public URL prefixes allowed without auth
	publicPrefixes = []string{
		"/v1/public/",
		"/v1/oauth/",
	}
	publicGRPCExact = map[string]bool{
		"/salonapp.v1.UserService/LoginUser":               true,
		"/salonapp.v1.UserService/RefreshToken":            true,
		"/salonapp.v1.UserService/VerifyEmailOTP":          true,
		"/salonapp.v1.UserService/ResendEmailVerification": true,
		"/salonapp.v1.UserService/CreateUser":              true,
		"/salonapp.v1.UserService/RecoverPassword":         true,
		"/salonapp.v1.UserService/ResetPassword":           true,
		"/salonapp.v1.UserService/LoginWithPhone":          true,
		"/salonapp.v1.UserService/RegisterPhoneUser":       true,
		"/salonapp.v1.UserService/RequestPhoneOTP":         true,
		"/salonapp.v1.UserService/VerifyPhoneOTP":          true,
		"/salonapp.v1.OAuthService/GetOAuthURL":            true,
	}
	publicGRPCPrefixes = []string{
		"/salonapp.v1.PublicService/",
		"/salonapp.v1.OAuthService/",
	}
	grpcRoleRules = map[string][]string{
		// "/salonapp.v1.UserService/GetUser": {string(entities.RoleSuperuser)},
	}
)

type AuthMiddleware struct {
	jwtRepository repositories.JWTRepository
	userRepo      repositories.UserRepository // interface to get user by ID
}

func NewAuthMiddleware(jwtRepository repositories.JWTRepository, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtRepository: jwtRepository,
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
			writeJSONError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		claims, err := m.jwtRepository.ValidateToken(token)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := m.userRepo.GetByID(r.Context(), claims.UserID)
		if err != nil || !user.IsActive {
			writeJSONError(w, http.StatusUnauthorized, "user not found or inactive")
			return
		}

		// Add user to context
		ctx := util.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// writeJSONError writes a standard JSON error body recognizable by the frontend.
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// GRPC interceptor for authentication
func (m *AuthMiddleware) GRPCAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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
	requiredRoles := RequiredGRPCRoles(info.FullMethod)
	if !util.HasRole(user, requiredRoles...) {
		return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
	}

	// Add user to context
	ctx = util.WithUser(ctx, user)
	return handler(ctx, req)
}

func RequiredGRPCRoles(method string) []string {
	return grpcRoleRules[method]
}

// GRPC interceptor for role-based authorization
func (m *AuthMiddleware) GRPCRoleInterceptor(requiredRoles ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		user := util.UserFromContext(ctx)
		if user == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if !util.HasRole(user, requiredRoles...) {
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
