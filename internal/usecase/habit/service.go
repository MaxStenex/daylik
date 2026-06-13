package habit

import (
	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
	habitLog "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type service struct {
	habitRepo    domain.Repository
	habitLogRepo habitLog.Repository
}

func New(habitRepo domain.Repository, habitLogRepo habitLog.Repository) *service {
	return &service{
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
	}
}
