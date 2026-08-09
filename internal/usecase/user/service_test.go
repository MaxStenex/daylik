package user

import (
	"testing"
	"time"

	mockrefreshtoken "github.com/maximrozinkevich/daylik/internal/generated/mocks/domain/refresh_token"
	mockuser "github.com/maximrozinkevich/daylik/internal/generated/mocks/domain/user"
	mockusecase "github.com/maximrozinkevich/daylik/internal/generated/mocks/usecase/user"
)

const testRefreshTTL = time.Hour

type serviceMocks struct {
	userRepo  *mockuser.MockRepository
	tokenRepo *mockrefreshtoken.MockRepository
	tokens    *mockusecase.MockTokenManager
	txm       *mockusecase.MockTxManager
}

func newTestService(t *testing.T) (*service, *serviceMocks) {
	m := &serviceMocks{
		userRepo:  mockuser.NewMockRepository(t),
		tokenRepo: mockrefreshtoken.NewMockRepository(t),
		tokens:    mockusecase.NewMockTokenManager(t),
		txm:       mockusecase.NewMockTxManager(t),
	}
	return New(m.userRepo, m.tokenRepo, m.tokens, m.txm, testRefreshTTL), m
}
