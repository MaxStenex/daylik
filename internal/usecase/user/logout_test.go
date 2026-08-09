package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	refresh_token "github.com/maximrozinkevich/daylik/internal/domain/refresh_token"
)

func TestService_Logout(t *testing.T) {
	userID := uuid.New()
	errRepo := errors.New("db down")

	tests := []struct {
		name    string
		in      LogoutInput
		setup   func(m *serviceMocks)
		wantErr error
	}{
		{
			name: "success",
			in:   LogoutInput{UserID: userID, RefreshToken: "raw-refresh"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().
					DeleteByHashAndUserID(mock.Anything, hashToken("raw-refresh"), userID).
					Return(nil)
			},
		},
		{
			name: "token not found is ignored",
			in:   LogoutInput{UserID: userID, RefreshToken: "raw-refresh"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().
					DeleteByHashAndUserID(mock.Anything, hashToken("raw-refresh"), userID).
					Return(refresh_token.ErrNotFound)
			},
		},
		{
			name: "repo error",
			in:   LogoutInput{UserID: userID, RefreshToken: "raw-refresh"},
			setup: func(m *serviceMocks) {
				m.tokenRepo.EXPECT().
					DeleteByHashAndUserID(mock.Anything, hashToken("raw-refresh"), userID).
					Return(errRepo)
			},
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, m := newTestService(t)
			tt.setup(m)

			err := srv.Logout(context.Background(), tt.in)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
