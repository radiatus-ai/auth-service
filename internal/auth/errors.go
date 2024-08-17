package auth

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorizedEmail  = errors.New("email not authorized")
	// Add other auth-related errors here
)
