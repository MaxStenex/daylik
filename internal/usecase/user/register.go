package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	domain "github.com/maximrozinkevich/daylik/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrEmailTaken       = errors.New("email already registered")
)

const minPasswordLen = 8

func (srv *service) Register(ctx context.Context, in RegisterInput) (RegisterOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	if _, err := mail.ParseAddress(email); err != nil {
		return RegisterOutput{}, ErrInvalidEmail
	}

	if len(in.Password) < minPasswordLen {
		return RegisterOutput{}, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("register: hash password: %w", err)
	}

	u := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err = srv.repo.Create(ctx, u); err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			return RegisterOutput{}, ErrEmailTaken
		}
		return RegisterOutput{}, fmt.Errorf("register: create user: %w", err)
	}

	return RegisterOutput{
		ID:    u.ID.String(),
		Email: u.Email,
	}, nil
}
