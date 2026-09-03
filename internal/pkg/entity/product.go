package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID     uint           `gorm:"index;not null" json:"store_id"`
	CategoryID  uint           `gorm:"index;not null" json:"category_id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(15,2);not null;index" json:"price"`
	Quantity    int            `gorm:"default:0;not null" json:"quantity"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url"`
	Store       *Store         `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ProductLogs []ProductLog   `gorm:"foreignKey:ProductID" json:"product_logs,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
