package habit_log

import (
	"context"
	"fmt"

	"github.com/maximrozinkevich/daylik/internal/helpers"
)

func (s *service) Update(ctx context.Context, in UpdateInput) error {
	log, err := s.habitLogRepo.FindByID(ctx, in.ID)
	if err != nil {
		return ErrNotFound
	}

	if log.UserID != in.UserID {
		return ErrForbidden
	}

	if !helpers.IsToday(log.CreatedAt) {
		return ErrUpdateOutdated
	}

	habit, err := s.habitRepo.FindByID(ctx, log.HabitID)
	if err != nil {
		return fmt.Errorf("habit_log: update: find habit: %w", err)
	}

	if in.CompletedCount <= 0 {
		return ErrInvalidCompletedCount
	}

	if in.CompletedCount > habit.DailyTarget {
		return ErrCompletedCountTooLarge
	}

	log.Unit = habit.Unit
	log.DailyTarget = habit.DailyTarget
	log.CompletedCount = in.CompletedCount
	if log.CompletedCount > log.DailyTarget {
		log.CompletedCount = log.DailyTarget
	}

	if err := s.habitLogRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("habit_log: update: %w", err)
	}

	return nil
}
