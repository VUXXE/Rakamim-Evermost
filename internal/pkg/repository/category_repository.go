package repository

import (
	"errors"

	"evermos-api/internal/pkg/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	FindByName(name string) (*entity.Category, error)
	FindByID(id uint) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]entity.Category, int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

func (r *categoryRepository) FindAll(limit, offset int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	query := r.db.Model(&entity.Category{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("id desc").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
