package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
)

func (srv *service) Refresh(ctx context.Context, in RefreshInput) (RefreshOutput, error) {
	rt, err := srv.tokenRepo.FindByHash(ctx, hashToken(in.RefreshToken))
	if err != nil {
		if errors.Is(err, refresh_token.ErrNotFound) {
			return RefreshOutput{}, ErrInvalidRefreshToken
		}
		return RefreshOutput{}, fmt.Errorf("refresh: find token: %w", err)
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = srv.tokenRepo.DeleteByHash(ctx, rt.Hash)
		return RefreshOutput{}, ErrInvalidRefreshToken
	}

	newRaw, err := srv.tokens.GenerateRefresh()
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("refresh: generate refresh token: %w", err)
	}

	newRT := &refresh_token.RefreshToken{
		UserID:    rt.UserID,
		Hash:      hashToken(newRaw),
		ExpiresAt: time.Now().Add(srv.refreshTTL),
	}

	if err = srv.txm.RunInTx(ctx, func(ctx context.Context) error {
		if err := srv.tokenRepo.DeleteByHash(ctx, rt.Hash); err != nil {
			return fmt.Errorf("delete old token: %w", err)
		}
		if err := srv.tokenRepo.Create(ctx, newRT); err != nil {
			return fmt.Errorf("store new token: %w", err)
		}
		return nil
	}); err != nil {
		return RefreshOutput{}, fmt.Errorf("refresh: rotate token: %w", err)
	}

	accessToken, err := srv.tokens.IssueAccess(rt.UserID)
	if err != nil {
		return RefreshOutput{}, fmt.Errorf("refresh: issue access token: %w", err)
	}

	return RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: newRaw,
	}, nil
}
