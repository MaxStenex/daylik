package user

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNotFound       = errors.New("user not found")
)
