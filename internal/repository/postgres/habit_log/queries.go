package habit_log

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit_log"
)

func (r *repository) Create(ctx context.Context, h *domain.HabitLog) (*domain.HabitLog, error) {
	rows, err := r.pool.Query(ctx,
		`INSERT INTO habits_log (user_id, habit_id, completed_count)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		h.UserID, h.HabitID, h.CompletedCount,
	)
	if err != nil {
		return nil, fmt.Errorf("habit_log repo: create: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createRow])
	if err != nil {
		return nil, fmt.Errorf("habit_log repo: create: %w", err)
	}

	h.ID = row.ID
	h.CreatedAt = row.CreatedAt
	return h, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.HabitLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, habit_id, completed_count, created_at
		 FROM habits_log WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("habit_log repo: find by id: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[habitLogRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("habit_log repo: find by id: %w", err)
	}

	return row.toDomain(), nil
}

func (r *repository) FindTodayByHabitIDs(ctx context.Context, userID uuid.UUID, habitIDs []uuid.UUID) (map[uuid.UUID]*domain.HabitLog, error) {
	logs := make(map[uuid.UUID]*domain.HabitLog, len(habitIDs))
	if len(habitIDs) == 0 {
		return logs, nil
	}

	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT ON (habit_id) id, user_id, habit_id, completed_count, created_at
		 FROM habits_log
		 WHERE user_id = $1 AND habit_id = ANY($2) AND created_at::date = CURRENT_DATE
		 ORDER BY habit_id, created_at DESC`,
		userID, habitIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("habit_log repo: find today by habit ids: %w", err)
	}

	logRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[habitLogRow])
	if err != nil {
		return nil, fmt.Errorf("habit_log repo: find today by habit ids: %w", err)
	}

	for _, row := range logRows {
		logs[row.HabitID] = row.toDomain()
	}
	return logs, nil
}

func (r *repository) Update(ctx context.Context, h *domain.HabitLog) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits_log SET completed_count = $1 WHERE id = $2`,
		h.CompletedCount, h.ID,
	)
	if err != nil {
		return fmt.Errorf("habit_log repo: update: %w", err)
	}

	return nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM habits_log WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("habit_log repo: delete: %w", err)
	}

	return nil
}
