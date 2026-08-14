package domain

import (
	"context"
	"errors"
	"time"
)

const MaxFailedLoginAttempts = 3
const AccountLockDuration = 15 * time.Minute

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid username or password")
	ErrAccountLocked     = errors.New("account is locked")
)

// User represents the users table.
type User struct {
	UserID              int64
	Username            string
	Password            string
	Token               string
	FailedLoginAttempts int
	LockedUntil         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// IsLocked reports whether the account is currently locked out.
func (u *User) IsLocked(now time.Time) bool {
	return u.LockedUntil != nil && u.LockedUntil.After(now)
}

// UserRepository defines persistence operations for users.
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	IncrementFailedAttempts(ctx context.Context, userID int64) (int, error)
	ResetFailedAttempts(ctx context.Context, userID int64) error
	LockAccount(ctx context.Context, userID int64, until time.Time) error
	UpdateToken(ctx context.Context, userID int64, token string) error
}
