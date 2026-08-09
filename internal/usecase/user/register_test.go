package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	user "github.com/maximrozinkevich/daylik/internal/domain/user"
)

func TestService_Register(t *testing.T) {
	newID := uuid.New()
	errRepo := errors.New("db down")

	tests := []struct {
		name    string
		in      RegisterInput
		setup   func(m *serviceMocks)
		want    RegisterOutput
		wantErr error
	}{
		{
			name: "success normalizes email",
			in:   RegisterInput{Email: "  John@Example.COM ", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(u *user.User) bool {
						return u.Email == "john@example.com" &&
							bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password123")) == nil
					})).
					Run(func(_ context.Context, u *user.User) { u.ID = newID }).
					Return(nil)
			},
			want: RegisterOutput{ID: newID.String(), Email: "john@example.com"},
		},
		{
			name:    "invalid email",
			in:      RegisterInput{Email: "not-an-email", Password: "password123"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "empty email",
			in:      RegisterInput{Email: "   ", Password: "password123"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "password too short",
			in:      RegisterInput{Email: "john@example.com", Password: "1234567"},
			wantErr: ErrPasswordTooShort,
		},
		{
			name: "email already taken",
			in:   RegisterInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(user.ErrDuplicateEmail)
			},
			wantErr: ErrEmailTaken,
		},
		{
			name: "repo error",
			in:   RegisterInput{Email: "john@example.com", Password: "password123"},
			setup: func(m *serviceMocks) {
				m.userRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errRepo)
			},
			wantErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, m := newTestService(t)
			if tt.setup != nil {
				tt.setup(m)
			}

			got, err := srv.Register(context.Background(), tt.in)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
