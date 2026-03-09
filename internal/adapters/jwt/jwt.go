package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type claims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

type adapter struct {
	secret    []byte
	accessTTL time.Duration
}

func New(secret string, accessTTL time.Duration) *adapter {
	return &adapter{
		secret:    []byte(secret),
		accessTTL: accessTTL,
	}
}
