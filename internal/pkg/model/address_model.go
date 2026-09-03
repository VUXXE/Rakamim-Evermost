package model

import "time"

type CreateAddressRequest struct {
	JudulAlamat   string `json:"judul_alamat" validate:"required"`
	PenerimaNama  string `json:"penerima_nama" validate:"required"`
	PenerimaPhone string `json:"penerima_phone" validate:"required"`
	Provinsi      string `json:"provinsi" validate:"required"`
	ProvinsiID    string `json:"provinsi_id" validate:"required"`
	Kabupaten     string `json:"kabupaten" validate:"required"`
	KabupatenID   string `json:"kabupaten_id" validate:"required"`
	Kecamatan     string `json:"kecamatan" validate:"required"`
	KecamatanID   string `json:"kecamatan_id" validate:"required"`
	Kelurahan     string `json:"kelurahan" validate:"required"`
	KelurahanID   string `json:"kelurahan_id" validate:"required"`
	DetailAlamat  string `json:"detail_alamat" validate:"required"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	JudulAlamat   string `json:"judul_alamat"`
	PenerimaNama  string `json:"penerima_nama"`
	PenerimaPhone string `json:"penerima_phone"`
	Provinsi      string `json:"provinsi"`
	ProvinsiID    string `json:"provinsi_id"`
	Kabupaten     string `json:"kabupaten"`
	KabupatenID   string `json:"kabupaten_id"`
	Kecamatan     string `json:"kecamatan"`
	KecamatanID   string `json:"kecamatan_id"`
	Kelurahan     string `json:"kelurahan"`
	KelurahanID   string `json:"kelurahan_id"`
	DetailAlamat  string `json:"detail_alamat"`
	IsDefault     *bool  `json:"is_default"`
}

type AddressResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	JudulAlamat   string    `json:"judul_alamat"`
	PenerimaNama  string    `json:"penerima_nama"`
	PenerimaPhone string    `json:"penerima_phone"`
	Provinsi      string    `json:"provinsi"`
	ProvinsiID    string    `json:"provinsi_id"`
	Kabupaten     string    `json:"kabupaten"`
	KabupatenID   string    `json:"kabupaten_id"`
	Kecamatan     string    `json:"kecamatan"`
	KecamatanID   string    `json:"kecamatan_id"`
	Kelurahan     string    `json:"kelurahan"`
	KelurahanID   string    `json:"kelurahan_id"`
	DetailAlamat  string    `json:"detail_alamat"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
