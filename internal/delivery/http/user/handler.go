package user

import (
	"context"
	"time"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type service interface {
	Register(ctx context.Context, in user.RegisterInput) (user.RegisterOutput, error)
	Login(ctx context.Context, in user.LoginInput) (user.LoginOutput, error)
	Refresh(ctx context.Context, in user.RefreshInput) (user.RefreshOutput, error)
	Logout(ctx context.Context, in user.LogoutInput) error
}

type Handler struct {
	srv          service
	refreshTTL   time.Duration
	cookieSecure bool
}

func New(srv service, refreshTTL time.Duration, cookieSecure bool) *Handler {
	return &Handler{srv: srv, refreshTTL: refreshTTL, cookieSecure: cookieSecure}
}
