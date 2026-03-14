package user

import (
	"crypto/sha256"
	"encoding/hex"
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

type service struct {
	userRepo   user.Repository
	tokenRepo  refresh_token.Repository
	tokens     TokenManager
	refreshTTL time.Duration
}

func New(
	userRepo user.Repository,
	tokenRepo refresh_token.Repository,
	tokens TokenManager,
	refreshTTL time.Duration,
) *service {
	return &service{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		tokens:     tokens,
		refreshTTL: refreshTTL,
	}
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}
