package util

import (
	"strings"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

// HasRole checks if user has any of the required roles
func HasRole(user *entities.User, requiredRoles ...string) bool {
	if len(requiredRoles) == 0 {
		return true
	}
	userRoles := make(map[string]bool)
	for _, role := range user.Roles {
		userRoles[strings.ToLower(role)] = true
	}

	if userRoles[string(entities.RoleSuperuser)] {
		return true
	}
	for _, required := range requiredRoles {
		if userRoles[strings.ToLower(required)] {
			return true
		}
	}

	return false
}

// HasAllRoles checks if user has all required roles
func HasAllRoles(user *entities.User, requiredRoles ...string) bool {
	userRoles := make(map[string]bool)
	for _, role := range user.Roles {
		userRoles[strings.ToLower(role)] = true
	}

	if userRoles[string(entities.RoleSuperuser)] {
		return true
	}
	for _, required := range requiredRoles {
		if !userRoles[strings.ToLower(required)] {
			return false
		}
	}

	return true
}
