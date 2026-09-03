package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryUseCase interface {
	CreateCategory(req model.CreateCategoryRequest) (*model.CategoryResponse, error)
	GetCategoryByID(id uint) (*model.CategoryResponse, error)
	GetAllCategories(limit, offset int) ([]model.CategoryResponse, int64, error)
	UpdateCategory(id uint, req model.UpdateCategoryRequest) (*model.CategoryResponse, error)
	DeleteCategory(id uint) error
}

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{categoryRepo: categoryRepo}
}

func (u *categoryUseCase) CreateCategory(req model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	existing, err := u.categoryRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category name already exists")
	}

	category := &entity.Category{
		Name: req.Name,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return toCategoryResponse(category), nil
}

func (u *categoryUseCase) GetCategoryByID(id uint) (*model.CategoryResponse, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	return toCategoryResponse(category), nil
}

func (u *categoryUseCase) GetAllCategories(limit, offset int) ([]model.CategoryResponse, int64, error) {
	categories, total, err := u.categoryRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.CategoryResponse, len(categories))
	for i, cat := range categories {
		result[i] = *toCategoryResponse(&cat)
	}

	return result, total, nil
}

func (u *categoryUseCase) UpdateCategory(id uint, req model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" && req.Name != category.Name {
		existing, err := u.categoryRepo.FindByName(req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("category name already exists")
		}
		category.Name = req.Name
	}

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return toCategoryResponse(category), nil
}

func (u *categoryUseCase) DeleteCategory(id uint) error {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return ErrCategoryNotFound
	}

	return u.categoryRepo.Delete(id)
}

func toCategoryResponse(c *entity.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
