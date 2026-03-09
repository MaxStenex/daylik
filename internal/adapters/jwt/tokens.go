package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (a *adapter) IssueAccess(userID uuid.UUID) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signed, err := t.SignedString(a.secret)
	if err != nil {
		return "", fmt.Errorf("token: sign access: %w", err)
	}

	return signed, nil
}

func (a *adapter) VerifyAccess(tokenStr string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.secret, nil
	})
	if err != nil || !t.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	c, ok := t.Claims.(*claims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	return c.UserID, nil
}

func (a *adapter) GenerateRefresh() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("token: generate refresh: %w", err)
	}
	return hex.EncodeToString(b), nil
}
