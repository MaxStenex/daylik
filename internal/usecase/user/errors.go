package user

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrEmailTaken          = errors.New("this email is already registered")
	ErrNotFound            = errors.New("user not found")
	ErrInternal            = errors.New("something went wrong")
)
