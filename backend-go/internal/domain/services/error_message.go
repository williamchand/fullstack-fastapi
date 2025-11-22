package services

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidRole       = errors.New("invalid role")
	ErrNoRolesProvided   = errors.New("no roles provided")
	ErrInvalidOAuthCode  = errors.New("invalid oauth code")
	ErrOAuthUnauthorized = errors.New("invalid oauth unauthorized")
)
