package habit

import (
	"context"
	"fmt"
)

func (s *service) List(ctx context.Context, in ListInput) (ListOutput, error) {
	habits, err := s.habitRepo.FindAllByUserID(ctx, in.UserID)
	if err != nil {
		return ListOutput{}, fmt.Errorf("habit: list: %w", err)
	}

	return ListOutput{Habits: habits}, nil
}
