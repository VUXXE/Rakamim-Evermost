package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
)

var (
	ErrStoreNotFound     = errors.New("store not found")
	ErrStoreAccessDenied = errors.New("forbidden: you do not own this store")
)

type StoreUseCase interface {
	CreateStore(userID uint, req model.CreateStoreRequest) (*model.StoreResponse, error)
	GetMyStore(userID uint) (*model.StoreResponse, error)
	UpdateStore(userID uint, storeID uint, req model.UpdateStoreRequest) (*model.StoreResponse, error)
	DeleteStore(userID uint, storeID uint) error
	GetAllStores(limit, offset int) ([]model.StoreResponse, int64, error)
	GetStoreByID(id uint) (*model.StoreResponse, error)
}

type storeUseCase struct {
	storeRepo repository.StoreRepository
}

func NewStoreUseCase(storeRepo repository.StoreRepository) StoreUseCase {
	return &storeUseCase{storeRepo: storeRepo}
}

func (u *storeUseCase) CreateStore(userID uint, req model.CreateStoreRequest) (*model.StoreResponse, error) {
	existing, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already has a store")
	}

	store := &entity.Store{
		UserID:  userID,
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	if err := u.storeRepo.Create(store); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &model.StoreResponse{
		ID:      store.ID,
		UserID:  store.UserID,
		Name:    store.Name,
		Address: store.Address,
		Phone:   store.Phone,
	}, nil
}

func (u *storeUseCase) GetMyStore(userID uint) (*model.StoreResponse, error) {
	store, err := u.storeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	return &model.StoreResponse{
		ID:      store.ID,
		UserID:  store.UserID,
		Name:    store.Name,
		Address: store.Address,
		Phone:   store.Phone,
	}, nil
}

func (u *storeUseCase) UpdateStore(userID uint, storeID uint, req model.UpdateStoreRequest) (*model.StoreResponse, error) {
	store, err := u.storeRepo.FindByID(storeID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	// Invariant: User can ONLY manage their own store
	if store.UserID != userID {
		return nil, ErrStoreAccessDenied
	}

	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Address != "" {
		store.Address = req.Address
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}

	if err := u.storeRepo.Update(store); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	return &model.StoreResponse{
		ID:      store.ID,
		UserID:  store.UserID,
		Name:    store.Name,
		Address: store.Address,
		Phone:   store.Phone,
	}, nil
}

func (u *storeUseCase) DeleteStore(userID uint, storeID uint) error {
	store, err := u.storeRepo.FindByID(storeID)
	if err != nil {
		return err
	}
	if store == nil {
		return ErrStoreNotFound
	}

	// Invariant: User can ONLY manage their own store
	if store.UserID != userID {
		return ErrStoreAccessDenied
	}

	return u.storeRepo.Delete(storeID)
}

func (u *storeUseCase) GetAllStores(limit, offset int) ([]model.StoreResponse, int64, error) {
	stores, total, err := u.storeRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.StoreResponse, len(stores))
	for i, st := range stores {
		result[i] = model.StoreResponse{
			ID:      st.ID,
			UserID:  st.UserID,
			Name:    st.Name,
			Address: st.Address,
			Phone:   st.Phone,
		}
	}

	return result, total, nil
}

func (u *storeUseCase) GetStoreByID(id uint) (*model.StoreResponse, error) {
	store, err := u.storeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	return &model.StoreResponse{
		ID:      store.ID,
		UserID:  store.UserID,
		Name:    store.Name,
		Address: store.Address,
		Phone:   store.Phone,
	}, nil
}
