package habit

import (
	"time"

	"github.com/google/uuid"
	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

type createRow struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type habitRow struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Name        string     `db:"name"`
	ExpReward   int64      `db:"exp_reward"`
	DailyTarget int64      `db:"daily_target"`
	Unit        string     `db:"unit"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (r habitRow) toDomain() *domain.Habit {
	return &domain.Habit{
		ID:          r.ID,
		UserID:      r.UserID,
		Name:        r.Name,
		ExpReward:   r.ExpReward,
		DailyTarget: r.DailyTarget,
		Unit:        r.Unit,
		CreatedAt:   r.CreatedAt,
		DeletedAt:   r.DeletedAt,
	}
}
