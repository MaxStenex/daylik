package refresh_token

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	domain "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	"github.com/maximrozinkevich/daylik/pkg/postgres"
)

func (r *repository) Create(ctx context.Context, t *domain.RefreshToken) error {
	_, err := postgres.Q(ctx, r.pool).Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, hash, expires_at)
		 VALUES ($1, $2, $3)`,
		t.UserID, t.Hash, t.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: create: %w", err)
	}
	return nil
}

func (r *repository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, hash, expires_at, created_at
		 FROM refresh_tokens WHERE hash = $1`,
		hash,
	)
	if err != nil {
		return nil, fmt.Errorf("refresh token repo: find by hash: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[refreshTokenRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("refresh token repo: find by hash: %w", err)
	}

	return row.toDomain(), nil
}

func (r *repository) DeleteByHash(ctx context.Context, hash string) error {
	_, err := postgres.Q(ctx, r.pool).Exec(ctx,
		`DELETE FROM refresh_tokens WHERE hash = $1`,
		hash,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: delete by hash: %w", err)
	}
	return nil
}

func (r *repository) DeleteByHashAndUserID(ctx context.Context, hash string, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE hash = $1 AND user_id = $2`,
		hash, userID,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: delete by hash and user id: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *repository) PruneOldest(ctx context.Context, userID uuid.UUID, maxCount int) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens
		 WHERE user_id = $1
		   AND id NOT IN (
		     SELECT id FROM refresh_tokens
		     WHERE user_id = $1
		     ORDER BY created_at DESC
		     LIMIT $2
		   )`,
		userID, maxCount,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: prune oldest: %w", err)
	}
	return nil
}

func (r *repository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: delete by user id: %w", err)
	}
	return nil
}
