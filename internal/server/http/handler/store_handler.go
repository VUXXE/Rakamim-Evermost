package handler

import (
	"errors"
	"strconv"

	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/usecase"
	"evermos-api/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	storeUseCase usecase.StoreUseCase
}

func NewStoreHandler(storeUseCase usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{storeUseCase: storeUseCase}
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	var req model.CreateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Name == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "name is required")
	}

	store, err := h.storeUseCase.CreateStore(userID, req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "failed to create store", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "store created successfully", store)
}

func (h *StoreHandler) GetMyStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	store, err := h.storeUseCase.GetMyStore(userID)
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "store not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch store", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", store)
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	storeID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid store id", err.Error())
	}

	var req model.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	store, err := h.storeUseCase.UpdateStore(userID, uint(storeID), req)
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "store not found", err.Error())
		}
		if errors.Is(err, usecase.ErrStoreAccessDenied) {
			return helper.Error(c, fiber.StatusForbidden, "forbidden", err.Error())
		}
		return helper.Error(c, fiber.StatusBadRequest, "failed to update store", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "store updated successfully", store)
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	storeID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid store id", err.Error())
	}

	err = h.storeUseCase.DeleteStore(userID, uint(storeID))
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "store not found", err.Error())
		}
		if errors.Is(err, usecase.ErrStoreAccessDenied) {
			return helper.Error(c, fiber.StatusForbidden, "forbidden", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to delete store", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "store deleted successfully", nil)
}

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	stores, total, err := h.storeUseCase.GetAllStores(limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch stores", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, stores)
}

func (h *StoreHandler) GetStoreByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid store id", err.Error())
	}

	store, err := h.storeUseCase.GetStoreByID(uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "store not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch store", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", store)
}
