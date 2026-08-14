package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/config"
	userhttp "api/internal/app/user/delivery/http"
	"api/internal/app/user/repository"
	"api/internal/app/user/usecase"
	"api/pkg/jwt"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	userRepo := repository.NewUserRepository(pool)
	userUsecase := usecase.NewUserUsecase(userRepo, func(userID int64, username string) (string, error) {
		return jwt.GenerateToken(cfg.JWTSecret, userID, username)
	})
	userHandler := userhttp.NewUserHandler(userUsecase)

	app := fiber.New()
	userHandler.RegisterRoutes(app)

	log.Printf("starting server on :%s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
