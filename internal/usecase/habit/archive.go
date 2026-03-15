package habit

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

func (s *service) Archive(ctx context.Context, in ArchiveInput) error {
	h, err := s.habitRepo.FindByID(ctx, in.HabitID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("habit: archive: find: %w", err)
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

	if err := s.habitRepo.ArchiveByID(ctx, in.HabitID); err != nil {
		return fmt.Errorf("habit: archive: %w", err)
	}

	return nil
}
