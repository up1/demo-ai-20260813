package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a domain.UserRepository backed by PostgreSQL.
func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = `
		SELECT user_id, username, password, token, failed_login_attempts, locked_until, created_at, updated_at
		FROM users
		WHERE username = $1`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&u.UserID,
		&u.Username,
		&u.Password,
		&u.Token,
		&u.FailedLoginAttempts,
		&u.LockedUntil,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}
	return &u, nil
}

func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID int64) (int, error) {
	const query = `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1, updated_at = now()
		WHERE user_id = $1
		RETURNING failed_login_attempts`

	var attempts int
	if err := r.pool.QueryRow(ctx, query, userID).Scan(&attempts); err != nil {
		return 0, fmt.Errorf("increment failed attempts: %w", err)
	}
	return attempts, nil
}

func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID int64) error {
	const query = `
		UPDATE users
		SET failed_login_attempts = 0, locked_until = NULL, updated_at = now()
		WHERE user_id = $1`

	if _, err := r.pool.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("reset failed attempts: %w", err)
	}
	return nil
}

func (r *userRepository) LockAccount(ctx context.Context, userID int64, until time.Time) error {
	const query = `
		UPDATE users
		SET locked_until = $2, updated_at = now()
		WHERE user_id = $1`

	if _, err := r.pool.Exec(ctx, query, userID, until); err != nil {
		return fmt.Errorf("lock account: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateToken(ctx context.Context, userID int64, token string) error {
	const query = `
		UPDATE users
		SET token = $2, updated_at = now()
		WHERE user_id = $1`

	if _, err := r.pool.Exec(ctx, query, userID, token); err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	return nil
}
