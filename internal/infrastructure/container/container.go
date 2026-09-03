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
	UserRepo    repository.UserRepository
	StoreRepo   repository.StoreRepository
	AddressRepo repository.AddressRepository

	// UseCases
	AuthUseCase    usecase.AuthUseCase
	UserUseCase    usecase.UserUseCase
	StoreUseCase   usecase.StoreUseCase
	AddressUseCase usecase.AddressUseCase

	// Handlers
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	StoreHandler   *handler.StoreHandler
	AddressHandler *handler.AddressHandler
}

func SetupContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(db, userRepo, storeRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	addressUseCase := usecase.NewAddressUseCase(db, addressRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	storeHandler := handler.NewStoreHandler(storeUseCase)
	addressHandler := handler.NewAddressHandler(addressUseCase)

	return &Container{
		DB:             db,
		UserRepo:       userRepo,
		StoreRepo:      storeRepo,
		AddressRepo:    addressRepo,
		AuthUseCase:    authUseCase,
		UserUseCase:    userUseCase,
		StoreUseCase:   storeUseCase,
		AddressUseCase: addressUseCase,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		StoreHandler:   storeHandler,
		AddressHandler: addressHandler,
	}
}
