package entity

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	Phone     string         `gorm:"type:varchar(32)" json:"phone"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Products  []Product      `gorm:"foreignKey:StoreID" json:"products,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
