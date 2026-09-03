package container

import (
	"evermos-api/internal/pkg/repository"
	"evermos-api/internal/pkg/usecase"
	"evermos-api/internal/server/http/handler"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepo  repository.UserRepository
	StoreRepo repository.StoreRepository

	// UseCases
	AuthUseCase usecase.AuthUseCase

	// Handlers
	AuthHandler *handler.AuthHandler
}

func SetupContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(db, userRepo, storeRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUseCase)

	return &Container{
		DB:          db,
		UserRepo:    userRepo,
		StoreRepo:   storeRepo,
		AuthUseCase: authUseCase,
		AuthHandler: authHandler,
	}
}
