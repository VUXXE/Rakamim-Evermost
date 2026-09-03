package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"evermos-api/internal/helper"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/usecase"
	"evermos-api/internal/server/http/middleware"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: productUseCase}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	var req model.CreateProductRequest

	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		catID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 32)
		price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
		qty, _ := strconv.Atoi(c.FormValue("quantity"))

		req.CategoryID = uint(catID)
		req.Name = c.FormValue("name")
		req.Description = c.FormValue("description")
		req.Price = price
		req.Quantity = qty

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			imagePath, err := saveProductImage(c, userID, file)
			if err != nil {
				return helper.Error(c, fiber.StatusBadRequest, "invalid image file", err.Error())
			}
			req.ImageURL = imagePath
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
		}
	}

	if req.CategoryID == 0 || req.Name == "" {
		return helper.Error(c, fiber.StatusBadRequest, "validation error", "category_id and name are required")
	}

	product, err := h.productUseCase.CreateProduct(userID, req)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "failed to create product", err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "product created successfully", product)
}

func (h *ProductHandler) GetMyProducts(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	products, total, err := h.productUseCase.GetMyProducts(userID, limit, offset)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch products", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, products)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	catID, _ := strconv.Atoi(c.Query("category_id", "0"))
	storeID, _ := strconv.Atoi(c.Query("store_id", "0"))

	filter := model.ProductFilter{
		Search:     c.Query("search"),
		CategoryID: uint(catID),
		StoreID:    uint(storeID),
		Sort:       c.Query("sort", "newest"),
		Limit:      limit,
		Offset:     offset,
	}

	products, total, err := h.productUseCase.GetAllProducts(filter)
	if err != nil {
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch products", err.Error())
	}

	return helper.Pagination(c, fiber.StatusOK, "success", total, limit, offset, products)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid product id", err.Error())
	}

	product, err := h.productUseCase.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "product not found", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to fetch product", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "success", product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid product id", err.Error())
	}

	var req model.UpdateProductRequest

	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if catVal := c.FormValue("category_id"); catVal != "" {
			catID, _ := strconv.ParseUint(catVal, 10, 32)
			cid := uint(catID)
			req.CategoryID = &cid
		}
		if name := c.FormValue("name"); name != "" {
			req.Name = name
		}
		if desc := c.FormValue("description"); desc != "" {
			req.Description = desc
		}
		if priceVal := c.FormValue("price"); priceVal != "" {
			price, _ := strconv.ParseFloat(priceVal, 64)
			req.Price = &price
		}
		if qtyVal := c.FormValue("quantity"); qtyVal != "" {
			qty, _ := strconv.Atoi(qtyVal)
			req.Quantity = &qty
		}

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			imagePath, err := saveProductImage(c, userID, file)
			if err != nil {
				return helper.Error(c, fiber.StatusBadRequest, "invalid image file", err.Error())
			}
			req.ImageURL = imagePath
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, fiber.StatusBadRequest, "invalid request body", err.Error())
		}
	}

	product, err := h.productUseCase.UpdateProduct(userID, uint(id), req)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "product not found", err.Error())
		}
		if errors.Is(err, usecase.ErrProductAccessDenied) {
			return helper.Error(c, fiber.StatusForbidden, "forbidden", err.Error())
		}
		return helper.Error(c, fiber.StatusBadRequest, "failed to update product", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return helper.Error(c, fiber.StatusUnauthorized, "unauthorized", "invalid token user")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, fiber.StatusBadRequest, "invalid product id", err.Error())
	}

	err = h.productUseCase.DeleteProduct(userID, uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return helper.Error(c, fiber.StatusNotFound, "product not found", err.Error())
		}
		if errors.Is(err, usecase.ErrProductAccessDenied) {
			return helper.Error(c, fiber.StatusForbidden, "forbidden", err.Error())
		}
		return helper.Error(c, fiber.StatusInternalServerError, "failed to delete product", err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "product deleted successfully", nil)
}

func saveProductImage(c *fiber.Ctx, userID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !validExts[ext] {
		return "", errors.New("unsupported file extension, allowed: .jpg, .jpeg, .png, .webp")
	}

	uploadDir := "./uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	cleanBase := filepath.Base(file.Filename)
	savedFileName := fmt.Sprintf("%d_%s", userID, cleanBase)
	dest := filepath.Join(uploadDir, savedFileName)

	if err := c.SaveFile(file, dest); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filepath.Join("uploads", "products", savedFileName), nil
}
