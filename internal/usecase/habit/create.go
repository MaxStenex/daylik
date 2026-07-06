package habit

import (
	"context"
	"strings"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
	"github.com/maximrozinkevich/daylik/pkg/logger"
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
		s.log.Error("habit: create", logger.Err(err))
		return CreateOutput{}, ErrInternal
	}

	return CreateOutput{Habit: *h}, nil
}
