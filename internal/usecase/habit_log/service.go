package habit_log

import (
	habit "github.com/maximrozinkevich/daylik/internal/domain/habit"
	habitLog "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type service struct {
	habitLogRepo habitLog.Repository
	habitRepo    habit.Repository
}

func New(
	habitLogRepo habitLog.Repository,
	habitRepo habit.Repository,
) *service {
	return &service{
		habitLogRepo: habitLogRepo,
		habitRepo:    habitRepo,
	}
}
