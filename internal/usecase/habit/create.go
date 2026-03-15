package habit

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

func (s *service) Create(ctx context.Context, in CreateInput) (CreateOutput, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return CreateOutput{}, ErrInvalidName
	}

	if in.ExpReward < 1 || in.ExpReward > 1000 {
		return CreateOutput{}, ErrInvalidExpReward
	}

	if in.DailyTarget <= 0 {
		return CreateOutput{}, ErrInvalidDailyTarget
	}

	unit := strings.TrimSpace(in.Unit)
	if unit == "" {
		return CreateOutput{}, ErrInvalidUnit
	}

	h := &domain.Habit{
		UserID:      in.UserID,
		Name:        name,
		ExpReward:   in.ExpReward,
		DailyTarget: in.DailyTarget,
		Unit:        unit,
	}

	if err := s.habitRepo.Create(ctx, h); err != nil {
		return CreateOutput{}, fmt.Errorf("habit: create: %w", err)
	}

	return CreateOutput{Habit: *h}, nil
}
