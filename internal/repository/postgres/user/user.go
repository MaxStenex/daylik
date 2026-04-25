package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	domain "github.com/maximrozinkevich/daylik/internal/domain/user"
)

func (r *repository) Create(ctx context.Context, u *domain.User) error {
	type createRow struct {
		ID        uuid.UUID `db:"id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	rows, err := r.pool.Query(ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id, created_at, updated_at`,
		u.Email, u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("user repo: create: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[createRow])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateEmail
		}
		return fmt.Errorf("user repo: create: %w", err)
	}

	u.ID = row.ID
	u.CreatedAt = row.CreatedAt
	u.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, email, password_hash, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("user repo: find by email: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("user repo: find by email: %w", err)
	}

	return row.toDomain(), nil
}
