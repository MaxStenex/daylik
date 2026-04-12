package refresh_token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	"github.com/maximrozinkevich/daylik/pkg/postgres"
)

type refreshTokenRow struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Hash      string    `db:"hash"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (r refreshTokenRow) toDomain() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        r.ID,
		UserID:    r.UserID,
		Hash:      r.Hash,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
	}
}

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, t *domain.RefreshToken) error {
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

func (r *Repository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
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

func (r *Repository) DeleteByHash(ctx context.Context, hash string) error {
	_, err := postgres.Q(ctx, r.pool).Exec(ctx,
		`DELETE FROM refresh_tokens WHERE hash = $1`,
		hash,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: delete by hash: %w", err)
	}
	return nil
}

func (r *Repository) DeleteByHashAndUserID(ctx context.Context, hash string, userID uuid.UUID) error {
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

func (r *Repository) PruneOldest(ctx context.Context, userID uuid.UUID, maxCount int) error {
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

func (r *Repository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("refresh token repo: delete by user id: %w", err)
	}
	return nil
}
