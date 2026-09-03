package usecase

import (
	"errors"
	"fmt"

	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
)

type UserUseCase interface {
	GetMe(userID uint) (*model.UserResponse, error)
	UpdateMe(userID uint, req model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteMe(userID uint) error
	GetAllUsers(limit, offset int) ([]model.UserResponse, int64, error)
	GetUserByID(id uint) (*model.UserResponse, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) GetMe(userID uint) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Phone:   user.Phone,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	}, nil
}

func (u *userUseCase) UpdateMe(userID uint, req model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// If email updated, ensure no collision with other users
	if req.Email != "" && req.Email != user.Email {
		existing, err := u.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}

	// If phone updated, ensure no collision with other users
	if req.Phone != "" && req.Phone != user.Phone {
		existing, err := u.userRepo.FindByPhone(req.Phone)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New("phone already taken")
		}
		user.Phone = req.Phone
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &model.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Phone:   user.Phone,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	}, nil
}

func (u *userUseCase) DeleteMe(userID uint) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return u.userRepo.Delete(userID)
}

func (u *userUseCase) GetAllUsers(limit, offset int) ([]model.UserResponse, int64, error) {
	users, total, err := u.userRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.UserResponse, len(users))
	for i, usr := range users {
		result[i] = model.UserResponse{
			ID:      usr.ID,
			Email:   usr.Email,
			Phone:   usr.Phone,
			Name:    usr.Name,
			IsAdmin: usr.IsAdmin,
		}
	}

	return result, total, nil
}

func (u *userUseCase) GetUserByID(id uint) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Phone:   user.Phone,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	}, nil
}
