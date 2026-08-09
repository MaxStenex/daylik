package habit

import (
	"context"
	"errors"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
	"github.com/maximrozinkevich/daylik/pkg/logger"
)

func (s *service) Delete(ctx context.Context, in DeleteInput) error {
	h, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ErrNotFound
		}
		s.log.Error("habit: delete: find", logger.Err(err))
		return ErrInternal
	}

	if h.DeletedAt != nil {
		return ErrNotFound
	}

	if h.UserID != in.UserID {
		return ErrForbidden
	}

	if err := s.habitRepo.DeleteByID(ctx, in.HabitID); err != nil {
		s.log.Error("habit: delete", logger.Err(err))
		return ErrInternal
	}

	return nil
}
