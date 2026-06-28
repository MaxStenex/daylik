package habit_log

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, habitLog *HabitLog) (*HabitLog, error)
	FindByID(ctx context.Context, id uuid.UUID) (*HabitLog, error)
	FindAllByHabitID(ctx context.Context, habitID uuid.UUID) ([]HabitLog, error)
	FindTodayByHabitIDs(ctx context.Context, userID uuid.UUID, habitIDs []uuid.UUID) (map[uuid.UUID]*HabitLog, error)
	Update(ctx context.Context, habitLog *HabitLog) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
