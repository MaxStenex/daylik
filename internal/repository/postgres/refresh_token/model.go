package refresh_token

import (
	"time"

	"github.com/google/uuid"
	domain "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
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
