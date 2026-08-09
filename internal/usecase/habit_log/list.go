package habit_log

import (
	"context"

	"github.com/maximrozinkevich/daylik/pkg/logger"
)

func (s *service) List(ctx context.Context, in ListInput) (ListOutput, error) {
	habit, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		return ListOutput{}, ErrInvalidHabit
	}

	if habit.DeletedAt != nil {
		return ListOutput{}, ErrInvalidHabit
	}

	if habit.UserID != in.UserID {
		return ListOutput{}, ErrForbidden
	}

	logs, err := s.habitLogRepo.FindAllByHabitID(ctx, in.HabitID)
	if err != nil {
		s.log.Error("habit_log: list", logger.Err(err))
		return ListOutput{}, ErrInternal
	}

	return ListOutput{HabitLogs: logs}, nil
}
