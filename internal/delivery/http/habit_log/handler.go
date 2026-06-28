package habits_log

import (
	"context"

	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit_log"
)

type service interface {
	Create(ctx context.Context, in usecase.CreateInput) (usecase.CreateOutput, error)
	List(ctx context.Context, in usecase.ListInput) (usecase.ListOutput, error)
	Update(ctx context.Context, in usecase.UpdateInput) error
}

type Handler struct {
	srv service
}

func New(srv service) *Handler {
	return &Handler{srv: srv}
}
