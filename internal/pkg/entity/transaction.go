package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	AddressID     uint           `gorm:"not null" json:"address_id"`
	InvoiceNumber string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoice_number"`
	TotalPrice    float64        `gorm:"type:decimal(15,2);not null" json:"total_price"`
	Status        string         `gorm:"type:varchar(32);default:'pending';not null" json:"status"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Address       *Address       `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	ProductLogs   []ProductLog   `gorm:"foreignKey:TransactionID" json:"product_logs,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
