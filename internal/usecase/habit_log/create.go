package habit_log

import (
	"context"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
	"github.com/maximrozinkevich/daylik/pkg/logger"
)

func (s *service) Create(ctx context.Context, in CreateInput) (CreateOutput, error) {
	habit, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		return CreateOutput{}, ErrInvalidHabit
	}

	if habit.DeletedAt != nil {
		return CreateOutput{}, ErrInvalidHabit
	}

	if habit.UserID != in.UserID {
		return CreateOutput{}, ErrForbidden
	}

	if in.CompletedCount <= 0 {
		return CreateOutput{}, ErrInvalidCompletedCount
	}

	if in.CompletedCount > habit.DailyTarget {
		return CreateOutput{}, ErrCompletedCountTooLarge
	}

	h := &domain.HabitLog{
		UserID:         in.UserID,
		HabitID:        in.HabitID,
		CompletedCount: in.CompletedCount,
		DailyTarget:    habit.DailyTarget,
		Unit:           habit.Unit,
	}

	created, err := s.habitLogRepo.Create(ctx, h)
	if err != nil {
		s.log.Error("habit_log: create", logger.Err(err))
		return CreateOutput{}, ErrInternal
	}

	return CreateOutput{HabitLog: *created}, nil
}
