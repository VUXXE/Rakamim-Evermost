package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductAccessDenied = errors.New("forbidden: you do not own this product")
)

type ProductUseCase interface {
	CreateProduct(userID uint, req model.CreateProductRequest) (*model.ProductResponse, error)
	GetMyProducts(userID uint, limit, offset int) ([]model.ProductResponse, int64, error)
	GetProductByID(id uint) (*model.ProductResponse, error)
	GetAllProducts(filter model.ProductFilter) ([]model.ProductResponse, int64, error)
	UpdateProduct(userID uint, productID uint, req model.UpdateProductRequest) (*model.ProductResponse, error)
	DeleteProduct(userID uint, productID uint) error
}

type productUseCase struct {
	productRepo  repository.ProductRepository
	storeRepo    repository.StoreRepository
	categoryRepo repository.CategoryRepository
}

func NewProductUseCase(
	productRepo repository.ProductRepository,
	storeRepo repository.StoreRepository,
	categoryRepo repository.CategoryRepository,
) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *productUseCase) CreateProduct(userID uint, req model.CreateProductRequest) (*model.ProductResponse, error) {
	store, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("user does not have an active store")
	}

	category, err := u.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if req.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	product := &entity.Product{
		StoreID:     store.ID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		ImageURL:    req.ImageURL,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	product.Store = store
	product.Category = category

	return toProductResponse(product), nil
}

func (u *productUseCase) GetMyProducts(userID uint, limit, offset int) ([]model.ProductResponse, int64, error) {
	store, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, err
	}
	if store == nil {
		return nil, 0, errors.New("user does not have an active store")
	}

	products, total, err := u.productRepo.FindByStoreID(store.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.ProductResponse, len(products))
	for i, p := range products {
		p.Store = store
		result[i] = *toProductResponse(&p)
	}

	return result, total, nil
}

func (u *productUseCase) GetProductByID(id uint) (*model.ProductResponse, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	return toProductResponse(product), nil
}

func (u *productUseCase) GetAllProducts(filter model.ProductFilter) ([]model.ProductResponse, int64, error) {
	products, total, err := u.productRepo.FindWithFilter(filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.ProductResponse, len(products))
	for i, p := range products {
		result[i] = *toProductResponse(&p)
	}

	return result, total, nil
}

func (u *productUseCase) UpdateProduct(userID uint, productID uint, req model.UpdateProductRequest) (*model.ProductResponse, error) {
	store, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errors.New("user does not have an active store")
	}

	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Zero-Trust Merchant Isolation: Verify product belongs to user's store
	if product.StoreID != store.ID {
		return nil, ErrProductAccessDenied
	}

	if req.CategoryID != nil && *req.CategoryID != product.CategoryID {
		category, err := u.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = *req.CategoryID
		product.Category = category
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return nil, errors.New("price cannot be negative")
		}
		product.Price = *req.Price
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return nil, errors.New("quantity cannot be negative")
		}
		product.Quantity = *req.Quantity
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return toProductResponse(product), nil
}

func (u *productUseCase) DeleteProduct(userID uint, productID uint) error {
	store, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return err
	}
	if store == nil {
		return errors.New("user does not have an active store")
	}

	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	// Zero-Trust Merchant Isolation
	if product.StoreID != store.ID {
		return ErrProductAccessDenied
	}

	return u.productRepo.Delete(productID)
}

func toProductResponse(p *entity.Product) *model.ProductResponse {
	resp := &model.ProductResponse{
		ID:          p.ID,
		StoreID:     p.StoreID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.Store != nil {
		resp.Store = &model.StoreResponse{
			ID:      p.Store.ID,
			UserID:  p.Store.UserID,
			Name:    p.Store.Name,
			Address: p.Store.Address,
			Phone:   p.Store.Phone,
		}
	}

	if p.Category != nil {
		resp.Category = &model.CategoryResponse{
			ID:        p.Category.ID,
			Name:      p.Category.Name,
			CreatedAt: p.Category.CreatedAt,
			UpdatedAt: p.Category.UpdatedAt,
		}
	}

	return resp
}
