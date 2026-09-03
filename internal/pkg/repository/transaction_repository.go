package repository

import (
	"errors"

	"evermos-api/internal/pkg/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithTx(tx *gorm.DB, txRecord *entity.Transaction) error
	FindByIDAndUserID(id uint, userID uint) (*entity.Transaction, error)
	FindByID(id uint) (*entity.Transaction, error)
	FindAllByUserID(userID uint, limit, offset int) ([]entity.Transaction, int64, error)
	FindAll(status string, limit, offset int) ([]entity.Transaction, int64, error)
	Update(txRecord *entity.Transaction) error
	UpdateWithTx(tx *gorm.DB, txRecord *entity.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateWithTx(tx *gorm.DB, txRecord *entity.Transaction) error {
	return tx.Create(txRecord).Error
}

func (r *transactionRepository) FindByIDAndUserID(id uint, userID uint) (*entity.Transaction, error) {
	var txRecord entity.Transaction
	err := r.db.Preload("Address").Preload("ProductLogs.Product").
		Where("id = ? AND user_id = ?", id, userID).
		First(&txRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &txRecord, nil
}

func (r *transactionRepository) FindByID(id uint) (*entity.Transaction, error) {
	var txRecord entity.Transaction
	err := r.db.Preload("Address").Preload("ProductLogs.Product").
		First(&txRecord, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &txRecord, nil
}

func (r *transactionRepository) FindAllByUserID(userID uint, limit, offset int) ([]entity.Transaction, int64, error) {
	var list []entity.Transaction
	var total int64

	query := r.db.Model(&entity.Transaction{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("Address").Preload("ProductLogs.Product").Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *transactionRepository) FindAll(status string, limit, offset int) ([]entity.Transaction, int64, error) {
	var list []entity.Transaction
	var total int64

	query := r.db.Model(&entity.Transaction{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("Address").Preload("ProductLogs.Product").Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *transactionRepository) Update(txRecord *entity.Transaction) error {
	return r.db.Save(txRecord).Error
}

func (r *transactionRepository) UpdateWithTx(tx *gorm.DB, txRecord *entity.Transaction) error {
	return tx.Save(txRecord).Error
}
