package handler

import (
	"strconv"

	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/usecase"
	"evermos-api/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	user, err := h.userUseCase.GetMe(userID)
	if err != nil {
		return helper.Error(c, fiber.StatusNotFound, "user not found", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", user)
}

func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	user, err := h.userUseCase.UpdateMe(userID, req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "failed to update user", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "user updated successfully", user)
}

func (h *UserHandler) DeleteMe(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	if err := h.userUseCase.DeleteMe(userID); err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to delete user", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "user deleted successfully", nil)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, total, err := h.userUseCase.GetAllUsers(limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch users", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, users)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid user id", err.Error())
	}

	user, err := h.userUseCase.GetUserByID(uint(id))
	if err != nil {
		return helper.Error(c, fiber.StatusNotFound, "user not found", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", user)
}
