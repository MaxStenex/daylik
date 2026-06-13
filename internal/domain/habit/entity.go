package habit

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type Habit struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	ExpReward   int64
	DailyTarget int64
	Unit        string
	TodayLog    *habit_log.HabitLog
	CreatedAt   time.Time
	DeletedAt   *time.Time
}
