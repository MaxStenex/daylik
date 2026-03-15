package habit

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, habit *Habit) error
	FindByID(ctx context.Context, id uuid.UUID) (*Habit, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Habit, error)
	Update(ctx context.Context, habit *Habit) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ArchiveByID(ctx context.Context, id uuid.UUID) error
}
