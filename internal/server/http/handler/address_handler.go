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

type AddressHandler struct {
	addressUseCase usecase.AddressUseCase
}

func NewAddressHandler(addressUseCase usecase.AddressUseCase) *AddressHandler {
	return &AddressHandler{addressUseCase: addressUseCase}
}

func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	var req model.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.JudulAlamat == "" || req.PenerimaNama == "" || req.PenerimaPhone == "" ||
		req.Provinsi == "" || req.ProvinsiID == "" || req.Kabupaten == "" || req.KabupatenID == "" ||
		req.Kecamatan == "" || req.KecamatanID == "" || req.Kelurahan == "" || req.KelurahanID == "" ||
		req.DetailAlamat == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "all address fields are required")
	}

	address, err := h.addressUseCase.CreateAddress(userID, req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "failed to create address", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "address created successfully", address)
}

func (h *AddressHandler) GetMyAddresses(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	addresses, total, err := h.addressUseCase.GetMyAddresses(userID, limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch addresses", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, addresses)
}

func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid address id", err.Error())
	}

	address, err := h.addressUseCase.GetAddressByID(userID, uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrAddressNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "address not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch address", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", address)
}

func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid address id", err.Error())
	}

	var req model.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	address, err := h.addressUseCase.UpdateAddress(userID, uint(id), req)
	if err != nil {
		if errors.Is(err, usecase.ErrAddressNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "address not found", err.Error())
		}
		return helper.Error(c, fiber.StatusBadRequest, "failed to update address", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "address updated successfully", address)
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid address id", err.Error())
	}

	err = h.addressUseCase.DeleteAddress(userID, uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrAddressNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "address not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to delete address", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "address deleted successfully", nil)
}
