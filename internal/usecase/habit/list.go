package habit

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *service) List(ctx context.Context, in ListInput) (ListOutput, error) {
	habits, err := s.habitRepo.FindAllByUserID(ctx, in.UserID)
	if err != nil {
		return ListOutput{}, fmt.Errorf("habit: list: %w", err)
	}

	habitIDs := make([]uuid.UUID, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}

	logs, err := s.habitLogRepo.FindTodayByHabitIDs(ctx, in.UserID, habitIDs)
	if err != nil {
		return ListOutput{}, fmt.Errorf("habit: list: %w", err)
	}

	for i := range habits {
		habits[i].TodayLog = logs[habits[i].ID]
	}

	return ListOutput{Habits: habits}, nil
}
