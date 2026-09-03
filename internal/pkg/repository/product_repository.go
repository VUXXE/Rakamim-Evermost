package repository

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	CreateWithTx(tx *gorm.DB, product *entity.Product) error
	FindByID(id uint) (*entity.Product, error)
	FindByIDWithLock(tx *gorm.DB, id uint) (*entity.Product, error)
	FindByStoreID(storeID uint, limit, offset int) ([]entity.Product, int64, error)
	FindWithFilter(filter model.ProductFilter) ([]entity.Product, int64, error)
	Update(product *entity.Product) error
	UpdateWithTx(tx *gorm.DB, product *entity.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) CreateWithTx(tx *gorm.DB, product *entity.Product) error {
	return tx.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Preload("Store").Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByIDWithLock(tx *gorm.DB, id uint) (*entity.Product, error) {
	var product entity.Product
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByStoreID(storeID uint, limit, offset int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.Model(&entity.Product{}).Where("store_id = ?", storeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("Category").Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindWithFilter(filter model.ProductFilter) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.Model(&entity.Product{})

	if filter.Search != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", filter.Search))
	}
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.StoreID > 0 {
		query = query.Where("store_id = ?", filter.StoreID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch filter.Sort {
	case "price_asc":
		query = query.Order("price asc")
	case "price_desc":
		query = query.Order("price desc")
	case "newest":
		query = query.Order("created_at desc")
	default:
		query = query.Order("created_at desc")
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Preload("Store").Preload("Category").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateWithTx(tx *gorm.DB, product *entity.Product) error {
	return tx.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Product{}, id).Error
}
