package habit_log

import (
	"time"

	"github.com/google/uuid"
)

type HabitLog struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	HabitID        uuid.UUID
	CompletedCount int64
	CreatedAt      time.Time
}
