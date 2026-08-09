package user

import (
	"context"
	"errors"
	"strings"
	"time"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	user "github.com/maximrozinkevich/daylik/internal/domain/user"
	"github.com/maximrozinkevich/daylik/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func (srv *service) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	u, err := srv.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return LoginOutput{}, ErrInvalidCredentials
		}
		srv.log.Error("login: find user", logger.Err(err))
		return LoginOutput{}, ErrInternal
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	accessToken, err := srv.tokens.IssueAccess(u.ID)
	if err != nil {
		srv.log.Error("login: issue access token", logger.Err(err))
		return LoginOutput{}, ErrInternal
	}

	rawRefresh, err := srv.tokens.GenerateRefresh()
	if err != nil {
		srv.log.Error("login: generate refresh token", logger.Err(err))
		return LoginOutput{}, ErrInternal
	}

	rt := &refresh_token.RefreshToken{
		UserID:    u.ID,
		Hash:      hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(srv.refreshTTL),
	}

	if err = srv.tokenRepo.Create(ctx, rt); err != nil {
		srv.log.Error("login: store refresh token", logger.Err(err))
		return LoginOutput{}, ErrInternal
	}

	_ = srv.tokenRepo.PruneOldest(ctx, u.ID, maxSessionsPerUser)

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
