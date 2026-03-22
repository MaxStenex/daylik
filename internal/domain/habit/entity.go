package habit

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	ExpReward   int64
	DailyTarget int64
	Unit        string
	CreatedAt   time.Time
	ArchivedAt  *time.Time
	DeletedAt   *time.Time
}
