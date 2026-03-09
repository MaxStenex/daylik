package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	user "github.com/maximrozinkevich/daylik/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
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

	u := &user.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err = srv.userRepo.Create(ctx, u); err != nil {
		if errors.Is(err, user.ErrDuplicateEmail) {
			return RegisterOutput{}, ErrEmailTaken
		}
		return RegisterOutput{}, fmt.Errorf("register: create user: %w", err)
	}

	return RegisterOutput{
		ID:    u.ID.String(),
		Email: u.Email,
	}, nil
}
