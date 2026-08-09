package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
	user "github.com/maximrozinkevich/daylik/internal/domain/user"
)

func TestService_Login(t *testing.T) {
	userID := uuid.New()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	require.NoError(t, err)
	existing := &user.User{ID: userID, Email: "john@example.com", PasswordHash: string(passwordHash)}
	errRepo := errors.New("db down")
	errIssueAccess := errors.New("no key")
	errGenerateRefresh := errors.New("entropy exhausted")

	tests := []struct {
		name    string
		in      LoginInput
		setup   func(m *serviceMocks)
		want    LoginOutput
		wantErr error
	}{
		{
			name: "success normalizes email",
			in:   LoginInput{Email: "  John@Example.COM ", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("access-token", nil)
				m.tokens.EXPECT().GenerateRefresh().Return("raw-refresh", nil)
				m.tokenRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(rt *refresh_token.RefreshToken) bool {
						return rt.UserID == userID &&
							rt.Hash == hashToken("raw-refresh") &&
							time.Until(rt.ExpiresAt) > testRefreshTTL-time.Minute
					})).
					Return(nil)
				m.tokenRepo.EXPECT().PruneOldest(mock.Anything, userID, maxSessionsPerUser).Return(nil)
			},
			want: LoginOutput{AccessToken: "access-token", RefreshToken: "raw-refresh"},
		},
		{
			name: "user not found",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(nil, user.ErrNotFound)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			in:   LoginInput{Email: "john@example.com", Password: "wrong-password"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "find user repo error",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(nil, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "issue access token error",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("", errIssueAccess)
			},
			wantErr: errIssueAccess,
		},
		{
			name: "generate refresh token error",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("access-token", nil)
				m.tokens.EXPECT().GenerateRefresh().Return("", errGenerateRefresh)
			},
			wantErr: errGenerateRefresh,
		},
		{
			name: "store refresh token error",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("access-token", nil)
				m.tokens.EXPECT().GenerateRefresh().Return("raw-refresh", nil)
				m.tokenRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(errRepo)
			},
			wantErr: errRepo,
		},
		{
			name: "prune error is ignored",
			in:   LoginInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().FindByEmail(mock.Anything, "john@example.com").Return(existing, nil)
				m.tokens.EXPECT().IssueAccess(userID).Return("access-token", nil)
				m.tokens.EXPECT().GenerateRefresh().Return("raw-refresh", nil)
				m.tokenRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
				m.tokenRepo.EXPECT().PruneOldest(mock.Anything, userID, maxSessionsPerUser).Return(errRepo)
			},
			want: LoginOutput{AccessToken: "access-token", RefreshToken: "raw-refresh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, m := newTestService(t)
			if tt.setup != nil {
				tt.setup(m)
			}

			got, err := srv.Login(context.Background(), tt.in)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
