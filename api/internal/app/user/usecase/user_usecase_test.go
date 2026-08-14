package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"api/internal/domain"
	"api/pkg/password"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) IncrementFailedAttempts(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockUserRepository) ResetFailedAttempts(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) LockAccount(ctx context.Context, userID int64, until time.Time) error {
	args := m.Called(ctx, userID, until)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func newHashedUser(id int64, username, plainPassword string) *domain.User {
	hashed, err := password.Hash(plainPassword)
	if err != nil {
		panic(err)
	}
	return &domain.User{
		UserID:   id,
		Username: username,
		Password: hashed,
		Address:  "123 Main St, City, Country",
	}
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepository)
	user := newHashedUser(1, "somkiat", "12345678")

	repo.On("GetByUsername", mock.Anything, "somkiat").Return(user, nil)
	repo.On("ResetFailedAttempts", mock.Anything, int64(1)).Return(nil)
	repo.On("UpdateToken", mock.Anything, int64(1), "signed-token").Return(nil)

	uc := NewUserUsecase(repo, func(userID int64, username string) (string, error) {
		return "signed-token", nil
	})

	result, err := uc.Login(context.Background(), "somkiat", "12345678")

	require.NoError(t, err)
	assert.Equal(t, int64(1), result.UserID)
	assert.Equal(t, "somkiat", result.Username)
	assert.Equal(t, "123 Main St, City, Country", result.Address)
	assert.Equal(t, "signed-token", result.Token)
	repo.AssertExpectations(t)
}

func TestLogin_UserNotFound_ReturnsInvalidCredential(t *testing.T) {
	repo := new(mockUserRepository)
	repo.On("GetByUsername", mock.Anything, "unknown").Return(nil, domain.ErrUserNotFound)

	uc := NewUserUsecase(repo, func(int64, string) (string, error) { return "", nil })

	result, err := uc.Login(context.Background(), "unknown", "12345678")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrInvalidCredential)
	repo.AssertExpectations(t)
}

func TestLogin_WrongPassword_ReturnsInvalidCredential(t *testing.T) {
	repo := new(mockUserRepository)
	user := newHashedUser(1, "somkiat", "12345678")

	repo.On("GetByUsername", mock.Anything, "somkiat").Return(user, nil)
	repo.On("IncrementFailedAttempts", mock.Anything, int64(1)).Return(1, nil)

	uc := NewUserUsecase(repo, func(int64, string) (string, error) { return "", nil })

	result, err := uc.Login(context.Background(), "somkiat", "wrongpassword")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrInvalidCredential)
	repo.AssertExpectations(t)
}

func TestLogin_AlreadyLocked_ReturnsAccountLocked(t *testing.T) {
	repo := new(mockUserRepository)
	user := newHashedUser(1, "somkiat", "12345678")
	lockedUntil := time.Now().Add(10 * time.Minute)
	user.LockedUntil = &lockedUntil

	repo.On("GetByUsername", mock.Anything, "somkiat").Return(user, nil)

	uc := NewUserUsecase(repo, func(int64, string) (string, error) { return "", nil }).(*userUsecase)
	uc.now = fixedNow(time.Now())

	result, err := uc.Login(context.Background(), "somkiat", "12345678")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrAccountLocked)
	repo.AssertExpectations(t)
}

func TestLogin_ThirdFailedAttempt_LocksAccount(t *testing.T) {
	repo := new(mockUserRepository)
	user := newHashedUser(1, "somkiat", "12345678")

	repo.On("GetByUsername", mock.Anything, "somkiat").Return(user, nil)
	repo.On("IncrementFailedAttempts", mock.Anything, int64(1)).Return(3, nil)
	repo.On("LockAccount", mock.Anything, int64(1), mock.AnythingOfType("time.Time")).Return(nil)

	uc := NewUserUsecase(repo, func(int64, string) (string, error) { return "", nil })

	result, err := uc.Login(context.Background(), "somkiat", "wrongpassword")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrAccountLocked)
	repo.AssertExpectations(t)
}

func TestLogin_RepositoryError_IsWrapped(t *testing.T) {
	repo := new(mockUserRepository)
	dbErr := errors.New("connection refused")
	repo.On("GetByUsername", mock.Anything, "somkiat").Return(nil, dbErr)

	uc := NewUserUsecase(repo, func(int64, string) (string, error) { return "", nil })

	result, err := uc.Login(context.Background(), "somkiat", "12345678")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.False(t, errors.Is(err, domain.ErrInvalidCredential))
	repo.AssertExpectations(t)
}
