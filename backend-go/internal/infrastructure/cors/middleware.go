package cors

import (
	"net/http"
	"strings"
)

// Middleware handles CORS headers
func Middleware(allowedOrigins string) func(http.Handler) http.Handler {
	origins := parseOrigins(allowedOrigins)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				if isOriginAllowed(origin, origins) {
					setCORSHeaders(w, origin)
					w.WriteHeader(http.StatusNoContent)
					return
				}
				http.Error(w, "CORS not allowed", http.StatusForbidden)
				return
			}

			// Handle actual requests
			if origin != "" && isOriginAllowed(origin, origins) {
				setCORSHeaders(w, origin)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")
}

func parseOrigins(allowedOrigins string) []string {
	if allowedOrigins == "*" {
		return []string{"*"}
	}
	return strings.Split(allowedOrigins, ",")
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}

	// Allow all origins
	if allowedOrigins[0] == "*" {
		return true
	}

	// Check if origin is in the allowed list
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}

	return false
}

