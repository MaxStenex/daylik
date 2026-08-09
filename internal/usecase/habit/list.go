package habit

import (
	"context"

	"github.com/google/uuid"

	"github.com/maximrozinkevich/daylik/pkg/logger"
)

func (s *service) List(ctx context.Context, in ListInput) (ListOutput, error) {
	habits, err := s.habitRepo.FindAllByUserID(ctx, in.UserID)
	if err != nil {
		s.log.Error("habit: list", logger.Err(err))
		return ListOutput{}, ErrInternal
	}

	habitIDs := make([]uuid.UUID, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}

	logs, err := s.habitLogRepo.FindTodayByHabitIDs(ctx, in.UserID, habitIDs)
	if err != nil {
		s.log.Error("habit: list", logger.Err(err))
		return ListOutput{}, ErrInternal
	}

	for i := range habits {
		habits[i].TodayLog = logs[habits[i].ID]
	}

	return ListOutput{Habits: habits}, nil
}
