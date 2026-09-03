package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:varchar(191);uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"phone"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin   bool           `gorm:"default:false" json:"is_admin"`
	Store     *Store         `gorm:"foreignKey:UserID" json:"store,omitempty"`
	Addresses []Address      `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
