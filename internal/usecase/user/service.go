package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	"github.com/maximrozinkevich/daylik/internal/domain/user"
)

const maxSessionsPerUser = 5

type TokenManager interface {
	IssueAccess(userID uuid.UUID) (string, error)
	GenerateRefresh() (string, error)
}

type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type service struct {
	userRepo   user.Repository
	tokenRepo  refresh_token.Repository
	tokens     TokenManager
	txm        TxManager
	refreshTTL time.Duration
	log        *slog.Logger
}

func New(
	userRepo user.Repository,
	tokenRepo refresh_token.Repository,
	tokens TokenManager,
	txm TxManager,
	refreshTTL time.Duration,
	log *slog.Logger,
) *service {
	return &service{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		tokens:     tokens,
		txm:        txm,
		refreshTTL: refreshTTL,
		log:        log,
	}
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}
