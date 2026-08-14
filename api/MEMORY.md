# API — Implementation Memory

## Feature: Login API (`req/login/feature_login_api.md`)

### What was built
- `POST /api/login` implemented with Clean Architecture: `cmd/main.go` → `config` → `internal/app/user/delivery/http` → `internal/app/user/usecase` → `internal/app/user/repository` → `internal/domain`.
- Stack: Go 1.26.5, Fiber v3.5.0, pgx/v5 + pgxpool, golang-jwt/jwt/v5, golang.org/x/crypto/bcrypt, testify, testcontainers-go (postgres module, `postgres:16-alpine`).

### Requirement gaps that needed a judgment call
- The spec's `users` table (user_id, username, password, token, created_at, updated_at) has no way to implement business rule #7 ("3 failed logins locks the account for 15 minutes"). Added two columns not in the original schema: `failed_login_attempts INTEGER` and `locked_until TIMESTAMPTZ`. See `migrations/0001_create_users_table.sql`. If the schema is later regenerated from the spec verbatim, this lockout logic will break — keep these columns or move the counter to a separate table/cache.
- The spec's validation message for a missing **username** is `"Please enter a valid email address."` even though the field is `username`, not email. Implemented verbatim as specified (not a bug) — do not "fix" this wording without checking with the requirement owner.
- No status code was specified for the lockout response; reused 400 with `status: "error"` to stay consistent with the other error responses in the spec.
- Token is a JWT (HS256, 24h expiry, `JWT_SECRET` env var) since the spec only shows an example JWT-shaped string but doesn't mandate a library/algorithm.

### Fiber v3 API notes (v3.5.0, differs from v2)
- Handlers take `fiber.Ctx` (interface), not `*fiber.Ctx`.
- Body binding: `c.Bind().Body(&req)`, not `c.BodyParser(&req)`.
- Context for downstream calls: `c.Context()` returns a standard `context.Context`.
- `app.Test(req, ...)` takes a variadic `fiber.TestConfig{Timeout, FailOnTimeout}`; default timeout is 1s if omitted — fine for local Fiber-in-process tests but be aware if a test does slow work per request.

### Testing setup
- Unit tests (`internal/app/user/usecase/user_usecase_test.go`) mock `domain.UserRepository` with `testify/mock`; the token generator is injected as a plain `func(userID int64, username string) (string, error)` so no JWT library is needed in unit tests.
- Integration tests (`integration_tests/user/`) spin up a real Postgres via testcontainers-go (`postgres.Run(ctx, "postgres:16-alpine", postgres.WithInitScripts(migrationPath), postgres.BasicWaitStrategies())`), then drive the real `fiber.App` end-to-end through `app.Test`. Requires Docker running locally; ~19s cold (image pull/container start) to run once.
- `postgres.WithInitScripts` runs the SQL migration file directly, so the integration tests and `migrations/0001_create_users_table.sql` must stay in sync — there is no separate migration runner in this project yet.
- Shared container/app setup lives in `setup_test.go`'s `TestMain`; each test calls `truncateUsers(t)` (+ `seedUser`) for isolation rather than restarting the container per test.

### Commands
- Unit tests only: `go test ./internal/...`
- Integration tests only (needs Docker): `go test ./integration_tests/...`
- Everything: `go test ./...`
