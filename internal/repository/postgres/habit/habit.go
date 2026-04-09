package habit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/maximrozinkevich/daylik/internal/domain/habit"
)

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

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, h *domain.Habit) error {
	rows, err := r.pool.Query(ctx,
		`INSERT INTO habits (user_id, name, exp_reward, daily_target, unit)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		h.UserID, h.Name, h.ExpReward, h.DailyTarget, h.Unit,
	)
	if err != nil {
		return fmt.Errorf("habit repo: create: %w", err)
	}

	type createRow struct {
		ID        uuid.UUID `db:"id"`
		CreatedAt time.Time `db:"created_at"`
	}
	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createRow])
	if err != nil {
		return fmt.Errorf("habit repo: create: %w", err)
	}

	h.ID = row.ID
	h.CreatedAt = row.CreatedAt
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Habit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, exp_reward, daily_target, unit, created_at, deleted_at
		 FROM habits WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("habit repo: find by id: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[habitRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("habit repo: find by id: %w", err)
	}

	return row.toDomain(), nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, exp_reward, daily_target, unit, created_at, deleted_at
		 FROM habits WHERE user_id = $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("habit repo: find all by user id: %w", err)
	}

	habitRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[habitRow])
	if err != nil {
		return nil, fmt.Errorf("habit repo: find all by user id: %w", err)
	}

	habits := make([]domain.Habit, len(habitRows))
	for i, row := range habitRows {
		habits[i] = *row.toDomain()
	}
	return habits, nil
}

func (r *Repository) Update(ctx context.Context, h *domain.Habit) error {
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

func (r *Repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits SET deleted_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("habit repo: delete: %w", err)
	}

	return nil
}
