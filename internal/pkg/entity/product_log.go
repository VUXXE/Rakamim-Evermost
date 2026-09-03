package entity

import (
	"time"
)

type ProductLog struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID uint      `gorm:"index;not null" json:"transaction_id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"type:decimal(15,2);not null" json:"price"`
	Product       *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
