package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api/internal/domain"
	"api/pkg/password"
)

// LoginResult is returned to the delivery layer on a successful login.
type LoginResult struct {
	UserID   int64
	Username string
	Token    string
}

// TokenGenerator produces a signed authentication token for a user.
type TokenGenerator func(userID int64, username string) (string, error)

// UserUsecase defines the login business logic.
type UserUsecase interface {
	Login(ctx context.Context, username, plainPassword string) (*LoginResult, error)
}

type userUsecase struct {
	repo          domain.UserRepository
	generateToken TokenGenerator
	now           func() time.Time
}

// NewUserUsecase wires the usecase with its repository and token generator.
func NewUserUsecase(repo domain.UserRepository, generateToken TokenGenerator) UserUsecase {
	return &userUsecase{
		repo:          repo,
		generateToken: generateToken,
		now:           time.Now,
	}
}

func (u *userUsecase) Login(ctx context.Context, username, plainPassword string) (*LoginResult, error) {
	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredential
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	now := u.now()
	if user.IsLocked(now) {
		return nil, domain.ErrAccountLocked
	}

	if !password.Matches(user.Password, plainPassword) {
		attempts, incErr := u.repo.IncrementFailedAttempts(ctx, user.UserID)
		if incErr != nil {
			return nil, fmt.Errorf("increment failed attempts: %w", incErr)
		}

		if attempts >= domain.MaxFailedLoginAttempts {
			if lockErr := u.repo.LockAccount(ctx, user.UserID, now.Add(domain.AccountLockDuration)); lockErr != nil {
				return nil, fmt.Errorf("lock account: %w", lockErr)
			}
			return nil, domain.ErrAccountLocked
		}

		return nil, domain.ErrInvalidCredential
	}

	if err := u.repo.ResetFailedAttempts(ctx, user.UserID); err != nil {
		return nil, fmt.Errorf("reset failed attempts: %w", err)
	}

	token, err := u.generateToken(user.UserID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	if err := u.repo.UpdateToken(ctx, user.UserID, token); err != nil {
		return nil, fmt.Errorf("update token: %w", err)
	}

	return &LoginResult{
		UserID:   user.UserID,
		Username: user.Username,
		Token:    token,
	}, nil
}
