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
	AuthUseCase  usecase.AuthUseCase
	UserUseCase  usecase.UserUseCase
	StoreUseCase usecase.StoreUseCase

	// Handlers
	AuthHandler  *handler.AuthHandler
	UserHandler  *handler.UserHandler
	StoreHandler *handler.StoreHandler
}

func SetupContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(db, userRepo, storeRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	storeHandler := handler.NewStoreHandler(storeUseCase)

	return &Container{
		DB:           db,
		UserRepo:     userRepo,
		StoreRepo:    storeRepo,
		AuthUseCase:  authUseCase,
		UserUseCase:  userUseCase,
		StoreUseCase: storeUseCase,
		AuthHandler:  authHandler,
		UserHandler:  userHandler,
		StoreHandler: storeHandler,
	}
}
