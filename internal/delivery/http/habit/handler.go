package habit

import (
	"context"

	usecase "github.com/maximrozinkevich/daylik/internal/usecase/habit"
)

type service interface {
	Create(ctx context.Context, in usecase.CreateInput) (usecase.CreateOutput, error)
	List(ctx context.Context, in usecase.ListInput) (usecase.ListOutput, error)
	Update(ctx context.Context, in usecase.UpdateInput) error
	Archive(ctx context.Context, in usecase.ArchiveInput) error
	Delete(ctx context.Context, in usecase.DeleteInput) error
}

type Handler struct {
	srv service
}

func NewHandler(srv service) *Handler {
	return &Handler{srv: srv}
}
