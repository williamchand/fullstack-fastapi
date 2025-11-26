package services

import "errors"

var (
	ErrUserNotActive           = errors.New("user is not active")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidEmailNotVerified = errors.New("invalid email not verified")
	ErrUserExists              = errors.New("user already exists")
	ErrInvalidRole             = errors.New("invalid role")
	ErrNoRolesProvided         = errors.New("no roles provided")
	ErrInvalidOAuthCode        = errors.New("invalid oauth code")
	ErrOAuthUnauthorized       = errors.New("invalid oauth unauthorized")
)
