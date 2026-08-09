package habit_log

import (
	"log/slog"

	habit "github.com/maximrozinkevich/daylik/internal/domain/habit"
	habitLog "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type service struct {
	habitLogRepo habitLog.Repository
	habitRepo    habit.Repository
	log          *slog.Logger
}

func New(
	habitLogRepo habitLog.Repository,
	habitRepo habit.Repository,
	log *slog.Logger,
) *service {
	return &service{
		habitLogRepo: habitLogRepo,
		habitRepo:    habitRepo,
		log:          log,
	}
}
