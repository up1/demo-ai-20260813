package user_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	userhttp "api/internal/app/user/delivery/http"
	"api/internal/app/user/repository"
	"api/internal/app/user/usecase"
	"api/pkg/jwt"
	"api/pkg/password"
)

const testJWTSecret = "test-secret"

var (
	testPool *pgxpool.Pool
	testApp  *fiber.App
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	ctx := context.Background()

	migrationPath, err := filepath.Abs("../../migrations/0001_create_users_table.sql")
	if err != nil {
		return 0, fmt.Errorf("resolve migration path: %w", err)
	}

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("app"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(migrationPath),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return 0, fmt.Errorf("start postgres container: %w", err)
	}
	defer func() { _ = pgContainer.Terminate(ctx) }()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return 0, fmt.Errorf("build connection string: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return 0, fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	testPool = pool

	userRepo := repository.NewUserRepository(pool)
	userUsecase := usecase.NewUserUsecase(userRepo, func(userID int64, username string) (string, error) {
		return jwt.GenerateToken(testJWTSecret, userID, username)
	})
	userHandler := userhttp.NewUserHandler(userUsecase)

	app := fiber.New()
	userHandler.RegisterRoutes(app)
	testApp = app

	return m.Run(), nil
}

// seedUser inserts a user with a bcrypt-hashed password and returns its id.
func seedUser(t *testing.T, username, plainPassword string) int64 {
	t.Helper()

	hashed, err := password.Hash(plainPassword)
	require.NoError(t, err)

	var id int64
	err = testPool.QueryRow(
		context.Background(),
		`INSERT INTO users (username, password, address) VALUES ($1, $2, $3) RETURNING user_id`,
		username, hashed, "123 Main St, City, Country",
	).Scan(&id)
	require.NoError(t, err)

	return id
}

// truncateUsers clears the users table so each test starts from a clean state.
func truncateUsers(t *testing.T) {
	t.Helper()

	_, err := testPool.Exec(context.Background(), `TRUNCATE TABLE users RESTART IDENTITY`)
	require.NoError(t, err)
}
