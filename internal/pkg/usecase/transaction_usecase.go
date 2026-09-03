package usecase

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
	"gorm.io/gorm"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionUseCase interface {
	CreateTransaction(userID uint, req model.CreateTransactionRequest) (*model.TransactionResponse, error)
	GetMyTransactions(userID uint, limit, offset int) ([]model.TransactionResponse, int64, error)
	GetMyTransactionByID(userID uint, id uint) (*model.TransactionResponse, error)
	UpdateTransactionStatus(userID uint, id uint, status string, isAdmin bool) (*model.TransactionResponse, error)
	GetAllTransactions(status string, limit, offset int) ([]model.TransactionResponse, int64, error)
}

type transactionUseCase struct {
	db          *gorm.DB
	txRepo      repository.TransactionRepository
	logRepo     repository.ProductLogRepository
	productRepo repository.ProductRepository
	addressRepo repository.AddressRepository
}

func NewTransactionUseCase(
	db *gorm.DB,
	txRepo repository.TransactionRepository,
	logRepo repository.ProductLogRepository,
	productRepo repository.ProductRepository,
	addressRepo repository.AddressRepository,
) TransactionUseCase {
	return &transactionUseCase{
		db:          db,
		txRepo:      txRepo,
		logRepo:     logRepo,
		productRepo: productRepo,
		addressRepo: addressRepo,
	}
}

func (u *transactionUseCase) CreateTransaction(userID uint, req model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	// 1. Verify shipping address belongs strictly to ordering user
	address, err := u.addressRepo.FindByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, errors.New("shipping address not found or does not belong to current user")
	}

	if len(req.Products) == 0 {
		return nil, errors.New("transaction must contain at least one product")
	}

	// 2. Sort items by ProductID ascending to eliminate deadlock risks during concurrent transactions
	items := make([]model.CheckoutItem, len(req.Products))
	copy(items, req.Products)
	sort.Slice(items, func(i, j int) bool {
		return items[i].ProductID < items[j].ProductID
	})

	var txRecord entity.Transaction
	var logs []entity.ProductLog
	var totalPrice float64

	// 3. Execute atomic database transaction
	err = u.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item.Quantity <= 0 {
				return errors.New("quantity must be greater than zero")
			}

			// Acquire pessimistic row lock
			product, err := u.productRepo.FindByIDWithLock(tx, item.ProductID)
			if err != nil {
				return err
			}
			if product == nil {
				return fmt.Errorf("product with ID %d not found", item.ProductID)
			}
			if product.Quantity < item.Quantity {
				return fmt.Errorf("insufficient stock for product '%s' (requested: %d, available: %d)", product.Name, item.Quantity, product.Quantity)
			}

			// Decrement inventory stock
			product.Quantity -= item.Quantity
			if err := u.productRepo.UpdateWithTx(tx, product); err != nil {
				return err
			}

			itemSubtotal := product.Price * float64(item.Quantity)
			totalPrice += itemSubtotal

			// Stage immutable product snapshot log
			logs = append(logs, entity.ProductLog{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				Price:     product.Price,
				Product:   product,
			})
		}

		invoiceNumber := fmt.Sprintf("INV-%d-%d", userID, time.Now().UnixNano())
		txRecord = entity.Transaction{
			UserID:        userID,
			AddressID:     req.AddressID,
			InvoiceNumber: invoiceNumber,
			TotalPrice:    totalPrice,
			Status:        "pending",
		}

		if err := u.txRepo.CreateWithTx(tx, &txRecord); err != nil {
			return err
		}

		for i := range logs {
			logs[i].TransactionID = txRecord.ID
		}

		if err := u.logRepo.CreateBatchWithTx(tx, logs); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	txRecord.Address = address
	txRecord.ProductLogs = logs

	return toTransactionResponse(&txRecord), nil
}

func (u *transactionUseCase) GetMyTransactions(userID uint, limit, offset int) ([]model.TransactionResponse, int64, error) {
	list, total, err := u.txRepo.FindAllByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.TransactionResponse, len(list))
	for i, item := range list {
		result[i] = *toTransactionResponse(&item)
	}

	return result, total, nil
}

func (u *transactionUseCase) GetMyTransactionByID(userID uint, id uint) (*model.TransactionResponse, error) {
	txRecord, err := u.txRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if txRecord == nil {
		return nil, ErrTransactionNotFound
	}

	return toTransactionResponse(txRecord), nil
}

func (u *transactionUseCase) UpdateTransactionStatus(userID uint, id uint, status string, isAdmin bool) (*model.TransactionResponse, error) {
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid transaction status, allowed: pending, completed, cancelled")
	}

	var txRecord *entity.Transaction
	var err error

	if isAdmin {
		txRecord, err = u.txRepo.FindByID(id)
	} else {
		txRecord, err = u.txRepo.FindByIDAndUserID(id, userID)
	}

	if err != nil {
		return nil, err
	}
	if txRecord == nil {
		return nil, ErrTransactionNotFound
	}

	txRecord.Status = status
	if err := u.txRepo.Update(txRecord); err != nil {
		return nil, err
	}

	return toTransactionResponse(txRecord), nil
}

func (u *transactionUseCase) GetAllTransactions(status string, limit, offset int) ([]model.TransactionResponse, int64, error) {
	list, total, err := u.txRepo.FindAll(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.TransactionResponse, len(list))
	for i, item := range list {
		result[i] = *toTransactionResponse(&item)
	}

	return result, total, nil
}

func toTransactionResponse(t *entity.Transaction) *model.TransactionResponse {
	resp := &model.TransactionResponse{
		ID:            t.ID,
		UserID:        t.UserID,
		AddressID:     t.AddressID,
		InvoiceNumber: t.InvoiceNumber,
		TotalPrice:    t.TotalPrice,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.Address != nil {
		resp.Address = toAddressResponse(t.Address)
	}

	if len(t.ProductLogs) > 0 {
		logs := make([]model.ProductLogResponse, len(t.ProductLogs))
		for i, l := range t.ProductLogs {
			logResp := model.ProductLogResponse{
				ID:            l.ID,
				TransactionID: l.TransactionID,
				ProductID:     l.ProductID,
				Quantity:      l.Quantity,
				Price:         l.Price,
				CreatedAt:     l.CreatedAt,
			}
			if l.Product != nil {
				logResp.Product = toProductResponse(l.Product)
			}
			logs[i] = logResp
		}
		resp.ProductLogs = logs
	}

	return resp
}
