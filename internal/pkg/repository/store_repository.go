package repository

import (
	"errors"

	"evermos-api/internal/pkg/entity"
	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(store *entity.Store) error
	CreateWithTx(tx *gorm.DB, store *entity.Store) error
	FindByUserID(userID uint) (*entity.Store, error)
	FindByID(id uint) (*entity.Store, error)
	Update(store *entity.Store) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]entity.Store, int64, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(store *entity.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepository) CreateWithTx(tx *gorm.DB, store *entity.Store) error {
	return tx.Create(store).Error
}

func (r *storeRepository) FindByUserID(userID uint) (*entity.Store, error) {
	var store entity.Store
	err := r.db.Where("user_id = ?", userID).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindByID(id uint) (*entity.Store, error) {
	var store entity.Store
	err := r.db.First(&store, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) Update(store *entity.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Store{}, id).Error
}

func (r *storeRepository) FindAll(limit, offset int) ([]entity.Store, int64, error) {
	var stores []entity.Store
	var total int64

	if err := r.db.Model(&entity.Store{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Model(&entity.Store{})
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}
