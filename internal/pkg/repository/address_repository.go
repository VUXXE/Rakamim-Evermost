package repository

import (
	"errors"

	"evermos-api/internal/pkg/entity"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *entity.Address) error
	CreateWithTx(tx *gorm.DB, address *entity.Address) error
	FindByIDAndUserID(id uint, userID uint) (*entity.Address, error)
	FindAllByUserID(userID uint, limit, offset int) ([]entity.Address, int64, error)
	FindByID(id uint) (*entity.Address, error)
	Update(address *entity.Address) error
	UpdateWithTx(tx *gorm.DB, address *entity.Address) error
	Delete(id uint, userID uint) error
	UnsetDefaultAddresses(tx *gorm.DB, userID uint) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *entity.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) CreateWithTx(tx *gorm.DB, address *entity.Address) error {
	return tx.Create(address).Error
}

func (r *addressRepository) FindByIDAndUserID(id uint, userID uint) (*entity.Address, error) {
	var address entity.Address
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) FindByID(id uint) (*entity.Address, error) {
	var address entity.Address
	err := r.db.First(&address, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) FindAllByUserID(userID uint, limit, offset int) ([]entity.Address, int64, error) {
	var addresses []entity.Address
	var total int64

	query := r.db.Model(&entity.Address{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("is_default desc, id desc").Find(&addresses).Error; err != nil {
		return nil, 0, err
	}

	return addresses, total, nil
}

func (r *addressRepository) Update(address *entity.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) UpdateWithTx(tx *gorm.DB, address *entity.Address) error {
	return tx.Save(address).Error
}

func (r *addressRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Address{}).Error
}

func (r *addressRepository) UnsetDefaultAddresses(tx *gorm.DB, userID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&entity.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}
