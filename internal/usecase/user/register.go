package user

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	user "github.com/maximrozinkevich/daylik/internal/domain/user"
	"github.com/maximrozinkevich/daylik/pkg/logger"
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
		srv.log.Error("register: hash password", logger.Err(err))
		return RegisterOutput{}, ErrInternal
	}

	u := &user.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err = srv.userRepo.Create(ctx, u); err != nil {
		if errors.Is(err, user.ErrDuplicateEmail) {
			return RegisterOutput{}, ErrEmailTaken
		}
		srv.log.Error("register: create user", logger.Err(err))
		return RegisterOutput{}, ErrInternal
	}

	return RegisterOutput{
		ID:    u.ID.String(),
		Email: u.Email,
	}, nil
}
