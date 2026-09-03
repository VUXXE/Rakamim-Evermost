package model

import "time"

type CheckoutItem struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type CreateTransactionRequest struct {
	AddressID uint           `json:"address_id" validate:"required"`
	Products  []CheckoutItem `json:"products" validate:"required,min=1"`
}

type UpdateTransactionStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type ProductLogResponse struct {
	ID            uint             `json:"id"`
	TransactionID uint             `json:"transaction_id"`
	ProductID     uint             `json:"product_id"`
	Quantity      int              `json:"quantity"`
	Price         float64          `json:"price"`
	Product       *ProductResponse `json:"product,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
}

type TransactionResponse struct {
	ID            uint                 `json:"id"`
	UserID        uint                 `json:"user_id"`
	AddressID     uint                 `json:"address_id"`
	InvoiceNumber string               `json:"invoice_number"`
	TotalPrice    float64              `json:"total_price"`
	Status        string               `json:"status"`
	Address       *AddressResponse     `json:"address,omitempty"`
	ProductLogs   []ProductLogResponse `json:"product_logs,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}
