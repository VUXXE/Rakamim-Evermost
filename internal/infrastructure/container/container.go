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
	UserRepo        repository.UserRepository
	StoreRepo       repository.StoreRepository
	AddressRepo     repository.AddressRepository
	CategoryRepo    repository.CategoryRepository
	ProductRepo     repository.ProductRepository
	TransactionRepo repository.TransactionRepository
	ProductLogRepo  repository.ProductLogRepository

	// UseCases
	AuthUseCase        usecase.AuthUseCase
	UserUseCase        usecase.UserUseCase
	StoreUseCase       usecase.StoreUseCase
	AddressUseCase     usecase.AddressUseCase
	CategoryUseCase    usecase.CategoryUseCase
	ProductUseCase     usecase.ProductUseCase
	TransactionUseCase usecase.TransactionUseCase
	RegionUseCase      usecase.RegionUseCase

	// Handlers
	AuthHandler        *handler.AuthHandler
	UserHandler        *handler.UserHandler
	StoreHandler       *handler.StoreHandler
	AddressHandler     *handler.AddressHandler
	CategoryHandler    *handler.CategoryHandler
	ProductHandler     *handler.ProductHandler
	TransactionHandler *handler.TransactionHandler
	SwaggerHandler     *handler.SwaggerHandler
	RegionHandler      *handler.RegionHandler
}

func SetupContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	productLogRepo := repository.NewProductLogRepository(db)

	// UseCases
	authUseCase := usecase.NewAuthUseCase(db, userRepo, storeRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	storeUseCase := usecase.NewStoreUseCase(storeRepo)
	addressUseCase := usecase.NewAddressUseCase(db, addressRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	productUseCase := usecase.NewProductUseCase(productRepo, storeRepo, categoryRepo)
	transactionUseCase := usecase.NewTransactionUseCase(db, transactionRepo, productLogRepo, productRepo, addressRepo)
	regionUseCase := usecase.NewRegionUseCase()

	// Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	storeHandler := handler.NewStoreHandler(storeUseCase)
	addressHandler := handler.NewAddressHandler(addressUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	swaggerHandler := handler.NewSwaggerHandler()
	regionHandler := handler.NewRegionHandler(regionUseCase)

	return &Container{
		DB:                 db,
		UserRepo:           userRepo,
		StoreRepo:          storeRepo,
		AddressRepo:        addressRepo,
		CategoryRepo:       categoryRepo,
		ProductRepo:        productRepo,
		TransactionRepo:    transactionRepo,
		ProductLogRepo:     productLogRepo,
		AuthUseCase:        authUseCase,
		UserUseCase:        userUseCase,
		StoreUseCase:       storeUseCase,
		AddressUseCase:     addressUseCase,
		CategoryUseCase:    categoryUseCase,
		ProductUseCase:     productUseCase,
		TransactionUseCase: transactionUseCase,
		RegionUseCase:      regionUseCase,
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		StoreHandler:       storeHandler,
		AddressHandler:     addressHandler,
		CategoryHandler:    categoryHandler,
		ProductHandler:     productHandler,
		TransactionHandler: transactionHandler,
		SwaggerHandler:     swaggerHandler,
		RegionHandler:      regionHandler,
	}
}
