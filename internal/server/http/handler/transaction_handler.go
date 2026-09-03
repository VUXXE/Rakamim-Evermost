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

type TransactionHandler struct {
	txUseCase usecase.TransactionUseCase
}

func NewTransactionHandler(txUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{txUseCase: txUseCase}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	var req model.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.AddressID == 0 || len(req.Products) == 0 {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "address_id and products are required")
	}

	txResp, err := h.txUseCase.CreateTransaction(userID, req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "transaction failed", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "transaction created successfully", txResp)
}

func (h *TransactionHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	list, total, err := h.txUseCase.GetMyTransactions(userID, limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch transactions", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, list)
}

func (h *TransactionHandler) GetMyTransactionByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid transaction id", err.Error())
	}

	txResp, err := h.txUseCase.GetMyTransactionByID(userID, uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrTransactionNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "transaction not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch transaction", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", txResp)
}

func (h *TransactionHandler) UpdateTransactionStatus(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid transaction id", err.Error())
	}

	var req model.UpdateTransactionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Status == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "status is required")
	}

	isAdmin := middleware.GetIsAdmin(c)

	txResp, err := h.txUseCase.UpdateTransactionStatus(userID, uint(id), req.Status, isAdmin)
	if err != nil {
		if errors.Is(err, usecase.ErrTransactionNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "transaction not found", err.Error())
		}
		return helper.Error(c, fiber.StatusBadRequest, "failed to update status", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "transaction status updated successfully", txResp)
}

func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	status := c.Query("status")

	list, total, err := h.txUseCase.GetAllTransactions(status, limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch transactions", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, list)
}
