package habit_log

import (
	"time"

	"github.com/google/uuid"
	domain "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

type createRow struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type habitLogRow struct {
	ID             uuid.UUID `db:"id"`
	UserID         uuid.UUID `db:"user_id"`
	HabitID        uuid.UUID `db:"habit_id"`
	Unit           string    `db:"unit"`
	CompletedCount int64     `db:"completed_count"`
	DailyTarget    int64     `db:"daily_target"`
	CreatedAt      time.Time `db:"created_at"`
}

func (r habitLogRow) toDomain() *domain.HabitLog {
	return &domain.HabitLog{
		ID:             r.ID,
		UserID:         r.UserID,
		HabitID:        r.HabitID,
		CompletedCount: r.CompletedCount,
		CreatedAt:      r.CreatedAt,
		Unit:           r.Unit,
		DailyTarget:    r.DailyTarget,
	}
}
