package repository

import (
	"evermos-api/internal/pkg/entity"
	"gorm.io/gorm"
)

type ProductLogRepository interface {
	CreateBatchWithTx(tx *gorm.DB, logs []entity.ProductLog) error
	FindByTransactionID(transactionID uint) ([]entity.ProductLog, error)
}

type productLogRepository struct {
	db *gorm.DB
}

func NewProductLogRepository(db *gorm.DB) ProductLogRepository {
	return &productLogRepository{db: db}
}

func (r *productLogRepository) CreateBatchWithTx(tx *gorm.DB, logs []entity.ProductLog) error {
	if len(logs) == 0 {
		return nil
	}
	return tx.Create(&logs).Error
}

func (r *productLogRepository) FindByTransactionID(transactionID uint) ([]entity.ProductLog, error) {
	var logs []entity.ProductLog
	err := r.db.Preload("Product").Where("transaction_id = ?", transactionID).Find(&logs).Error
	return logs, err
}
