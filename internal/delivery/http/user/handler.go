package user

import (
	"context"

	"github.com/maximrozinkevich/daylik/internal/usecase/user"
)

type service interface {
	Register(ctx context.Context, in user.RegisterInput) (user.RegisterOutput, error)
	Login(ctx context.Context, in user.LoginInput) (user.LoginOutput, error)
	Refresh(ctx context.Context, in user.RefreshInput) (user.RefreshOutput, error)
	Logout(ctx context.Context, in user.LogoutInput) error
}

type Handler struct {
	srv service
}

func NewHandler(srv service) *Handler {
	return &Handler{srv: srv}
}
