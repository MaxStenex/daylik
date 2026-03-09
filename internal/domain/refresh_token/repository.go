package refresh_token

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByHash(ctx context.Context, hash string) error
	DeleteByHashAndUserID(ctx context.Context, hash string, userID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	PruneOldest(ctx context.Context, userID uuid.UUID, maxCount int) error
}
