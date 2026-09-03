package handler

import (
	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Name == "" || req.Email == "" || req.Phone == "" || req.Password == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "name, email, phone, and password are required")
	}

	resp, err := h.authUseCase.Register(req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "registration failed", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "register success", resp)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Email == "" || req.Password == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "email and password are required")
	}

	resp, err := h.authUseCase.Login(req)
	if err != nil {
		return helper.Error(c, fiber.StatusUnauthorized, "login failed", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "login success", resp)
}
