package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	JudulAlamat   string         `gorm:"type:varchar(100);not null" json:"judul_alamat"`
	PenerimaNama  string         `gorm:"type:varchar(255);not null" json:"penerima_nama"`
	PenerimaPhone string         `gorm:"type:varchar(32);not null" json:"penerima_phone"`
	Provinsi      string         `gorm:"type:varchar(100);not null" json:"provinsi"`
	ProvinsiID    string         `gorm:"type:varchar(32);not null" json:"provinsi_id"`
	Kabupaten     string         `gorm:"type:varchar(100);not null" json:"kabupaten"`
	KabupatenID   string         `gorm:"type:varchar(32);not null" json:"kabupaten_id"`
	Kecamatan     string         `gorm:"type:varchar(100);not null" json:"kecamatan"`
	KecamatanID   string         `gorm:"type:varchar(32);not null" json:"kecamatan_id"`
	Kelurahan     string         `gorm:"type:varchar(100);not null" json:"kelurahan"`
	KelurahanID   string         `gorm:"type:varchar(32);not null" json:"kelurahan_id"`
	DetailAlamat  string         `gorm:"type:text;not null" json:"detail_alamat"`
	IsDefault     bool           `gorm:"default:false" json:"is_default"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
