package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	user "github.com/maximrozinkevich/daylik/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

func (srv *service) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	u, err := srv.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return LoginOutput{}, ErrInvalidCredentials
		}
		return LoginOutput{}, fmt.Errorf("login: find user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	accessToken, err := srv.tokens.IssueAccess(u.ID)
	if err != nil {
		return LoginOutput{}, fmt.Errorf("login: issue access token: %w", err)
	}

	rawRefresh, err := srv.tokens.GenerateRefresh()
	if err != nil {
		return LoginOutput{}, fmt.Errorf("login: generate refresh token: %w", err)
	}

	rt := &refresh_token.RefreshToken{
		UserID:    u.ID,
		Hash:      hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(srv.refreshTTL),
	}

	if err = srv.tokenRepo.Create(ctx, rt); err != nil {
		return LoginOutput{}, fmt.Errorf("login: store refresh token: %w", err)
	}

	_ = srv.tokenRepo.PruneOldest(ctx, u.ID, maxSessionsPerUser)

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
