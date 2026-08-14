package http

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"api/internal/app/user/usecase"
	"api/internal/domain"
	"api/pkg/response"
)

const minPasswordLength = 8

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Address  string `json:"address"`
	Token    string `json:"token"`
}

// UserHandler exposes HTTP handlers for the user feature.
type UserHandler struct {
	usecase usecase.UserUsecase
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

// RegisterRoutes wires the user routes onto the given router.
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/api/login", h.Login)
}

// Login handles POST /api/login.
func (h *UserHandler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid request body"))
	}

	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Please enter a valid email address."))
	}
	if len(req.Password) < minPasswordLength {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Password must be at least 8 characters long."))
	}

	result, err := h.usecase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredential):
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid username or password"))
		case errors.Is(err, domain.ErrAccountLocked):
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Account is locked due to multiple failed login attempts. Please try again in 15 minutes."))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("Login successful", loginData{
		UserID:   result.UserID,
		Username: result.Username,
		Address:  result.Address,
		Token:    result.Token,
	}))
}
