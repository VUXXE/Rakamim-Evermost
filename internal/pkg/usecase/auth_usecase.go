package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/entity"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
	"evermos-api/internal/utils"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Register(req model.RegisterRequest) (*model.AuthResponse, error)
	Login(req model.LoginRequest) (*model.AuthResponse, error)
}

type authUseCase struct {
	db        *gorm.DB
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
}

func NewAuthUseCase(db *gorm.DB, userRepo repository.UserRepository, storeRepo repository.StoreRepository) AuthUseCase {
	return &authUseCase{
		db:        db,
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
}

func (u *authUseCase) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	// Validate email uniqueness
	existingEmail, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already registered")
	}

	// Validate phone uniqueness
	existingPhone, err := u.userRepo.FindByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if existingPhone != nil {
		return nil, errors.New("phone already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		IsAdmin:  false,
	}

	var store *entity.Store

	// Atomic transaction: Create User & Auto-provision Store
	err = u.db.Transaction(func(tx *gorm.DB) error {
		if err := u.userRepo.CreateWithTx(tx, user); err != nil {
			return err
		}

		store = &entity.Store{
			UserID:  user.ID,
			Name:    fmt.Sprintf("%s's Store", user.Name),
			Address: "",
			Phone:   user.Phone,
		}

		if err := u.storeRepo.CreateWithTx(tx, store); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("registration transaction failed: %w", err)
	}

	// Generate JWT Token
	token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:      user.ID,
			Email:   user.Email,
			Phone:   user.Phone,
			Name:    user.Name,
			IsAdmin: user.IsAdmin,
		},
		Store: &model.StoreResponse{
			ID:      store.ID,
			UserID:  store.UserID,
			Name:    store.Name,
			Address: store.Address,
			Phone:   store.Phone,
		},
	}, nil
}

func (u *authUseCase) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	store, err := u.storeRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	var storeResp *model.StoreResponse
	if store != nil {
		storeResp = &model.StoreResponse{
			ID:      store.ID,
			UserID:  store.UserID,
			Name:    store.Name,
			Address: store.Address,
			Phone:   store.Phone,
		}
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			ID:      user.ID,
			Email:   user.Email,
			Phone:   user.Phone,
			Name:    user.Name,
			IsAdmin: user.IsAdmin,
		},
		Store: storeResp,
	}, nil
}
