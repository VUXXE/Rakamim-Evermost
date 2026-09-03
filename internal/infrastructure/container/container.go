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
	UserRepo     repository.UserRepository
	StoreRepo    repository.StoreRepository
	AddressRepo  repository.AddressRepository
	CategoryRepo repository.CategoryRepository

	// UseCases
	AuthUseCase     usecase.AuthUseCase
	UserUseCase     usecase.UserUseCase
	StoreUseCase    usecase.StoreUseCase
	AddressUseCase  usecase.AddressUseCase
	CategoryUseCase usecase.CategoryUseCase

	// Handlers
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	StoreHandler    *handler.StoreHandler
	AddressHandler  *handler.AddressHandler
	CategoryHandler *handler.CategoryHandler
}

func SetupContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(db, userRepo, storeRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	addressUseCase := usecase.NewAddressUseCase(db, addressRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	storeHandler := handler.NewStoreHandler(storeUseCase)
	addressHandler := handler.NewAddressHandler(addressUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)

	return &Container{
		DB:              db,
		UserRepo:        userRepo,
		StoreRepo:       storeRepo,
		AddressRepo:     addressRepo,
		CategoryRepo:    categoryRepo,
		AuthUseCase:     authUseCase,
		UserUseCase:     userUseCase,
		StoreUseCase:    storeUseCase,
		AddressUseCase:  addressUseCase,
		CategoryUseCase: categoryUseCase,
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		StoreHandler:    storeHandler,
		AddressHandler:  addressHandler,
		CategoryHandler: categoryHandler,
	}
}
