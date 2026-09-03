package handler

import (
	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/usecase"
	"github.com/gofiber/fiber/v2"
)

type RegionHandler struct {
	useCase usecase.RegionUseCase
}

func NewRegionHandler(u usecase.RegionUseCase) *RegionHandler {
	return &RegionHandler{useCase: u}
}

func (h *RegionHandler) GetProvinces(c *fiber.Ctx) error {
	ctx := c.UserContext()
	provinces, err := h.useCase.GetProvinces(ctx)
	if err != nil {
		return helper.Error(c, fiber.StatusBadGateway, "failed to fetch provinces", err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "provinces fetched successfully", provinces)
}

func (h *RegionHandler) GetRegencies(c *fiber.Ctx) error {
	ctx := c.UserContext()
	provinceID := c.Params("province_id")
	if provinceID == "" {
		return helper.Error(c, fiber.StatusBadRequest, "province_id is required", "missing path parameter")
	}

	regencies, err := h.useCase.GetRegencies(ctx, provinceID)
	if err != nil {
		return helper.Error(c, fiber.StatusBadGateway, "failed to fetch regencies", err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "regencies fetched successfully", regencies)
}

func (h *RegionHandler) GetDistricts(c *fiber.Ctx) error {
	ctx := c.UserContext()
	regencyID := c.Params("regency_id")
	if regencyID == "" {
		return helper.Error(c, fiber.StatusBadRequest, "regency_id is required", "missing path parameter")
	}

	districts, err := h.useCase.GetDistricts(ctx, regencyID)
	if err != nil {
		return helper.Error(c, fiber.StatusBadGateway, "failed to fetch districts", err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "districts fetched successfully", districts)
}

func (h *RegionHandler) GetVillages(c *fiber.Ctx) error {
	ctx := c.UserContext()
	districtID := c.Params("district_id")
	if districtID == "" {
		return helper.Error(c, fiber.StatusBadRequest, "district_id is required", "missing path parameter")
	}

	villages, err := h.useCase.GetVillages(ctx, districtID)
	if err != nil {
		return helper.Error(c, fiber.StatusBadGateway, "failed to fetch villages", err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "villages fetched successfully", villages)
}
