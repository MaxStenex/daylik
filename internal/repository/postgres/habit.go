package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

type HabitRepository struct {
	pool *pgxpool.Pool
}

func NewHabitRepository(pool *pgxpool.Pool) *HabitRepository {
	return &HabitRepository{pool: pool}
}

func (r *HabitRepository) Create(ctx context.Context, h *domain.Habit) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO habits (user_id, name, exp_reward, daily_target, unit)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		h.UserID, h.Name, h.ExpReward, h.DailyTarget, h.Unit,
	)

	if err := row.Scan(&h.ID, &h.CreatedAt); err != nil {
		return fmt.Errorf("habit repo: create: %w", err)
	}

	return nil
}

func (r *HabitRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Habit, error) {
	var h domain.Habit

	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, exp_reward, daily_target, unit, created_at, archived_at, deleted_at
		 FROM habits WHERE id = $1`,
		id,
	)

	if err := row.Scan(
		&h.ID, &h.UserID, &h.Name,
		&h.ExpReward, &h.DailyTarget, &h.Unit,
		&h.CreatedAt, &h.ArchivedAt, &h.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("habit repo: find by id: %w", err)
	}

	return &h, nil
}

func (r *HabitRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, exp_reward, daily_target, unit, created_at, archived_at, deleted_at
		 FROM habits WHERE user_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("habit repo: find all by user id: %w", err)
	}
	defer rows.Close()

	habits := make([]domain.Habit, 0)
	for rows.Next() {
		var h domain.Habit
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.Name,
			&h.ExpReward, &h.DailyTarget, &h.Unit,
			&h.CreatedAt, &h.ArchivedAt, &h.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("habit repo: scan: %w", err)
		}
		habits = append(habits, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("habit repo: rows: %w", err)
	}

	return habits, nil
}

func (r *HabitRepository) Update(ctx context.Context, h *domain.Habit) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits
		 SET name = $1, exp_reward = $2, daily_target = $3, unit = $4
		 WHERE id = $5`,
		h.Name, h.ExpReward, h.DailyTarget, h.Unit, h.ID,
	)
	if err != nil {
		return fmt.Errorf("habit repo: update: %w", err)
	}

	return nil
}

func (r *HabitRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits SET deleted_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("habit repo: delete: %w", err)
	}

	return nil
}

func (r *HabitRepository) ArchiveByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits SET archived_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("habit repo: archive: %w", err)
	}

	return nil
}
