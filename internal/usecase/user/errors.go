package user

import "errors"

var (
	ErrInvalidCredentials  = errors.New("Invalid email or password")
	ErrInvalidRefreshToken = errors.New("Invalid or expired refresh token")
	ErrInvalidEmail        = errors.New("Invalid email address")
	ErrPasswordTooShort    = errors.New("Password must be at least 8 characters")
	ErrEmailTaken          = errors.New("This email is already registered")
	ErrNotFound            = errors.New("User not found")
)
