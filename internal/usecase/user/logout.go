package user

import (
	"context"
	"errors"
	"fmt"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
)

func (srv *service) Logout(ctx context.Context, in LogoutInput) error {
	err := srv.tokenRepo.DeleteByHashAndUserID(ctx, hashToken(in.RefreshToken), in.UserID)
	if err != nil && !errors.Is(err, refresh_token.ErrNotFound) {
		return fmt.Errorf("logout: delete refresh token: %w", err)
	}
	return nil
}
