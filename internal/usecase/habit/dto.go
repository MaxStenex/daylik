package habit

import (
	"github.com/google/uuid"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

type CreateInput struct {
	UserID      uuid.UUID
	Name        string
	ExpReward   int64
	DailyTarget float64
	Unit        string
}

type CreateOutput struct {
	Habit domain.Habit
}

type ListInput struct {
	UserID uuid.UUID
}

type ListOutput struct {
	Habits []domain.Habit
}

type UpdateInput struct {
	HabitID     uuid.UUID
	UserID      uuid.UUID
	Name        string
	ExpReward   int64
	DailyTarget float64
	Unit        string
}

type ArchiveInput struct {
	HabitID uuid.UUID
	UserID  uuid.UUID
}

type DeleteInput struct {
	HabitID uuid.UUID
	UserID  uuid.UUID
}
