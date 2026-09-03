package handler

import (
	"errors"
	"strconv"

	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/usecase"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req model.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Name == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "category name is required")
	}

	category, err := h.categoryUseCase.CreateCategory(req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "failed to create category", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "category created successfully", category)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid category id", err.Error())
	}

	category, err := h.categoryUseCase.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "category not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch category", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", category)
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	categories, total, err := h.categoryUseCase.GetAllCategories(limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch categories", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, categories)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid category id", err.Error())
	}

	var req model.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Name == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "category name is required")
	}

	category, err := h.categoryUseCase.UpdateCategory(uint(id), req)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "category not found", err.Error())
		}
		return helper.Error(c, fiber.StatusBadRequest, "failed to update category", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid category id", err.Error())
	}

	err = h.categoryUseCase.DeleteCategory(uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "category not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to delete category", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "category deleted successfully", nil)
}
