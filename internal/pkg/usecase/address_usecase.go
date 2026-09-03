package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
	"gorm.io/gorm"
)

var (
	ErrAddressNotFound = errors.New("address not found")
)

type AddressUseCase interface {
	CreateAddress(userID uint, req model.CreateAddressRequest) (*model.AddressResponse, error)
	GetMyAddresses(userID uint, limit, offset int) ([]model.AddressResponse, int64, error)
	GetAddressByID(userID uint, id uint) (*model.AddressResponse, error)
	UpdateAddress(userID uint, id uint, req model.UpdateAddressRequest) (*model.AddressResponse, error)
	DeleteAddress(userID uint, id uint) error
}

type addressUseCase struct {
	db          *gorm.DB
	addressRepo repository.AddressRepository
}

func NewAddressUseCase(db *gorm.DB, addressRepo repository.AddressRepository) AddressUseCase {
	return &addressUseCase{
		db:          db,
		addressRepo: addressRepo,
	}
}

func (u *addressUseCase) CreateAddress(userID uint, req model.CreateAddressRequest) (*model.AddressResponse, error) {
	address := &entity.Address{
		UserID:        userID,
		JudulAlamat:   req.JudulAlamat,
		PenerimaNama:  req.PenerimaNama,
		PenerimaPhone: req.PenerimaPhone,
		Provinsi:      req.Provinsi,
		ProvinsiID:    req.ProvinsiID,
		Kabupaten:     req.Kabupaten,
		KabupatenID:   req.KabupatenID,
		Kecamatan:     req.Kecamatan,
		KecamatanID:   req.KecamatanID,
		Kelurahan:     req.Kelurahan,
		KelurahanID:   req.KelurahanID,
		DetailAlamat:  req.DetailAlamat,
		IsDefault:     req.IsDefault,
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault {
			if err := u.addressRepo.UnsetDefaultAddresses(tx, userID); err != nil {
				return err
			}
		}
		return u.addressRepo.CreateWithTx(tx, address)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return toAddressResponse(address), nil
}

func (u *addressUseCase) GetMyAddresses(userID uint, limit, offset int) ([]model.AddressResponse, int64, error) {
	addresses, total, err := u.addressRepo.FindAllByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.AddressResponse, len(addresses))
	for i, addr := range addresses {
		result[i] = *toAddressResponse(&addr)
	}

	return result, total, nil
}

func (u *addressUseCase) GetAddressByID(userID uint, id uint) (*model.AddressResponse, error) {
	address, err := u.addressRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, ErrAddressNotFound
	}

	return toAddressResponse(address), nil
}

func (u *addressUseCase) UpdateAddress(userID uint, id uint, req model.UpdateAddressRequest) (*model.AddressResponse, error) {
	address, err := u.addressRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}
	if address == nil {
		return nil, ErrAddressNotFound
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault != nil && *req.IsDefault {
			if err := u.addressRepo.UnsetDefaultAddresses(tx, userID); err != nil {
				return err
			}
			address.IsDefault = true
		} else if req.IsDefault != nil {
			address.IsDefault = false
		}

		if req.JudulAlamat != "" {
			address.JudulAlamat = req.JudulAlamat
		}
		if req.PenerimaNama != "" {
			address.PenerimaNama = req.PenerimaNama
		}
		if req.PenerimaPhone != "" {
			address.PenerimaPhone = req.PenerimaPhone
		}
		if req.Provinsi != "" {
			address.Provinsi = req.Provinsi
		}
		if req.ProvinsiID != "" {
			address.ProvinsiID = req.ProvinsiID
		}
		if req.Kabupaten != "" {
			address.Kabupaten = req.Kabupaten
		}
		if req.KabupatenID != "" {
			address.KabupatenID = req.KabupatenID
		}
		if req.Kecamatan != "" {
			address.Kecamatan = req.Kecamatan
		}
		if req.KecamatanID != "" {
			address.KecamatanID = req.KecamatanID
		}
		if req.Kelurahan != "" {
			address.Kelurahan = req.Kelurahan
		}
		if req.KelurahanID != "" {
			address.KelurahanID = req.KelurahanID
		}
		if req.DetailAlamat != "" {
			address.DetailAlamat = req.DetailAlamat
		}

		return u.addressRepo.UpdateWithTx(tx, address)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return toAddressResponse(address), nil
}

func (u *addressUseCase) DeleteAddress(userID uint, id uint) error {
	address, err := u.addressRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}
	if address == nil {
		return ErrAddressNotFound
	}

	return u.addressRepo.Delete(id, userID)
}

func toAddressResponse(a *entity.Address) *model.AddressResponse {
	return &model.AddressResponse{
		ID:            a.ID,
		UserID:        a.UserID,
		JudulAlamat:   a.JudulAlamat,
		PenerimaNama:  a.PenerimaNama,
		PenerimaPhone: a.PenerimaPhone,
		Provinsi:      a.Provinsi,
		ProvinsiID:    a.ProvinsiID,
		Kabupaten:     a.Kabupaten,
		KabupatenID:   a.KabupatenID,
		Kecamatan:     a.Kecamatan,
		KecamatanID:   a.KecamatanID,
		Kelurahan:     a.Kelurahan,
		KelurahanID:   a.KelurahanID,
		DetailAlamat:  a.DetailAlamat,
		IsDefault:     a.IsDefault,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}
