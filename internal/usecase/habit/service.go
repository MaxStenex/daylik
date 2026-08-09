package habit

import (
	"log/slog"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
	habitLog "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type service struct {
	habitRepo    domain.Repository
	habitLogRepo habitLog.Repository
	log          *slog.Logger
}

func New(habitRepo domain.Repository, habitLogRepo habitLog.Repository, log *slog.Logger) *service {
	return &service{
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
		log:          log,
	}
}
