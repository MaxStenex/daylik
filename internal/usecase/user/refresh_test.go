package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
)

func TestService_Refresh(t *testing.T) {
	userID := uuid.New()
	oldHash := hashToken("old-raw")
	errRepo := errors.New("db down")
	errIssueAccess := errors.New("no key")
	errGenerateRefresh := errors.New("entropy exhausted")

	validToken := func() *refresh_token.RefreshToken {
		return &refresh_token.RefreshToken{
			UserID:    userID,
			Hash:      oldHash,
			ExpiresAt: time.Now().Add(time.Hour),
		}
	}

	runTx := func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}

	tests := []struct {
		name    string
		in      RefreshInput
		setup   func(m *serviceMocks)
		want    RefreshOutput
		wantErr error
	}{
		{
			name: "success rotates token",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(validToken(), nil)
				m.tokens.EXPECT().GenerateRefresh().Return("new-raw", nil)
				m.txm.EXPECT().RunInTx(mock.Anything, mock.Anything).RunAndReturn(runTx)
				m.tokenRepo.EXPECT().DeleteByHash(mock.Anything, oldHash).Return(nil)
				m.tokenRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(rt *refresh_token.RefreshToken) bool {
						return rt.UserID == userID &&
							rt.Hash == hashToken("new-raw") &&
							time.Until(rt.ExpiresAt) > testRefreshTTL-time.Minute
					})).
					Return(nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("access-token", nil)
			},
			want: RefreshOutput{AccessToken: "access-token", RefreshToken: "new-raw"},
		},
		{
			name: "token not found",
			in:   RefreshInput{RefreshToken: "unknown-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().
					FindByHash(mock.Anything, hashToken("unknown-raw")).
					Return(nil, refresh_token.ErrNotFound)
			},
			wantErr: ErrInvalidRefreshToken,
		},
		{
			name: "find token repo error",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "expired token is deleted",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				expired := validToken()
				expired.ExpiresAt = time.Now().Add(-time.Minute)
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(expired, nil)
				m.tokenRepo.EXPECT().DeleteByHash(mock.Anything, oldHash).Return(nil)
			},
			wantErr: ErrInvalidRefreshToken,
		},
		{
			name: "generate refresh token error",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(validToken(), nil)
				m.tokens.EXPECT().GenerateRefresh().Return("", errGenerateRefresh)
			},
			wantErr: errGenerateRefresh,
		},
		{
			name: "rotation fails on delete",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(validToken(), nil)
				m.tokens.EXPECT().GenerateRefresh().Return("new-raw", nil)
				m.txm.EXPECT().RunInTx(mock.Anything, mock.Anything).RunAndReturn(runTx)
				m.tokenRepo.EXPECT().DeleteByHash(mock.Anything, oldHash).Return(errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "rotation fails on create",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(validToken(), nil)
				m.tokens.EXPECT().GenerateRefresh().Return("new-raw", nil)
				m.txm.EXPECT().RunInTx(mock.Anything, mock.Anything).RunAndReturn(runTx)
				m.tokenRepo.EXPECT().DeleteByHash(mock.Anything, oldHash).Return(nil)
				m.tokenRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "issue access token error",
			in:   RefreshInput{RefreshToken: "old-raw"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().FindByHash(mock.Anything, oldHash).Return(validToken(), nil)
				m.tokens.EXPECT().GenerateRefresh().Return("new-raw", nil)
				m.txm.EXPECT().RunInTx(mock.Anything, mock.Anything).RunAndReturn(runTx)
				m.tokenRepo.EXPECT().DeleteByHash(mock.Anything, oldHash).Return(nil)
				m.tokenRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("", errIssueAccess)
			},
			wantErr: errIssueAccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, m := newTestService(t)
			if tt.setup != nil {
				tt.setup(m)
			}

			got, err := srv.Refresh(context.Background(), tt.in)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
