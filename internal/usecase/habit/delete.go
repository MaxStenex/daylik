package habit

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

func (s *service) Delete(ctx context.Context, in DeleteInput) error {
	h, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("habit: delete: find: %w", err)
	}

	if h.DeletedAt != nil {
		return ErrNotFound
	}

	if h.UserID != in.UserID {
		return ErrForbidden
	}

	if err := s.habitRepo.DeleteByID(ctx, in.HabitID); err != nil {
		return fmt.Errorf("habit: delete: %w", err)
	}

	return nil
}
