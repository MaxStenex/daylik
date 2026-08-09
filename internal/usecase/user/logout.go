package user

import (
	"context"
	"errors"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	"github.com/maximrozinkevich/daylik/pkg/logger"
)

func (srv *service) Logout(ctx context.Context, in LogoutInput) error {
	err := srv.tokenRepo.DeleteByHashAndUserID(ctx, hashToken(in.RefreshToken), in.UserID)
	if err != nil && !errors.Is(err, refresh_token.ErrNotFound) {
		srv.log.Error("logout: delete refresh token", logger.Err(err))
		return ErrInternal
	}
	return nil
}
