package habit

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

func (s *service) Update(ctx context.Context, in UpdateInput) error {
	h, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("habit: update: find: %w", err)
	}

	if h.DeletedAt != nil {
		return ErrNotFound
	}

	if h.UserID != in.UserID {
		return ErrForbidden
	}

	if h.ArchivedAt != nil {
		return ErrAlreadyArchived
	}

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return ErrInvalidName
	}

	if in.ExpReward < 1 || in.ExpReward > 1000 {
		return ErrInvalidExpReward
	}

	if in.DailyTarget <= 0 {
		return ErrInvalidDailyTarget
	}

	unit := strings.TrimSpace(in.Unit)
	if unit == "" {
		return ErrInvalidUnit
	}

	h.Name = name
	h.ExpReward = in.ExpReward
	h.DailyTarget = in.DailyTarget
	h.Unit = unit

	if err := s.habitRepo.Update(ctx, h); err != nil {
		return fmt.Errorf("habit: update: %w", err)
	}

	return nil
}
