package habit_log

import (
	"github.com/google/uuid"
	domain "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type CreateInput struct {
	UserID         uuid.UUID
	HabitID        uuid.UUID
	CompletedCount int64
}

type CreateOutput struct {
	HabitLog domain.HabitLog
}

type UpdateInput struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CompletedCount int64
}

type ListInput struct {
	UserID  uuid.UUID
	HabitID uuid.UUID
}

type ListOutput struct {
	HabitLogs []domain.HabitLog
}
