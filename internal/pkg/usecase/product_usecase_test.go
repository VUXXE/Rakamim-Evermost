package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
	"evermos-api/internal/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductUseCase_FullLifecycle_AndIsolation(t *testing.T) {
	ctx := setupAllTestDB(t)
	categoryRepo := repository.NewCategoryRepository(ctx.db)
	productRepo := repository.NewProductRepository(ctx.db)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	productUC := usecase.NewProductUseCase(productRepo, ctx.storeRepo, categoryRepo)

	ts := time.Now().UnixNano()
	// Create Category
	cat, err := categoryUC.CreateCategory(model.CreateCategoryRequest{
		Name: fmt.Sprintf("Electronics_%d", ts%100000),
	})
	require.NoError(t, err)

	// User A
	userA, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Seller A",
		Email:    fmt.Sprintf("seller_a_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0897%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// User B
	userB, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Seller B",
		Email:    fmt.Sprintf("seller_b_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0898%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// 1. User A creates product
	prodA, err := productUC.CreateProduct(userA.User.ID, model.CreateProductRequest{
		CategoryID:  cat.ID,
		Name:        "Gaming Laptop Pro",
		Description: "High-end laptop",
		Price:       25000000,
		Quantity:    10,
		ImageURL:    "uploads/products/laptop.jpg",
	})
	require.NoError(t, err)
	assert.Equal(t, "Gaming Laptop Pro", prodA.Name)
	assert.Equal(t, userA.Store.ID, prodA.StoreID)

	// 2. User B tries to update User A's product -> MUST FAIL with ErrProductAccessDenied
	newPrice := float64(1000)
	_, err = productUC.UpdateProduct(userB.User.ID, prodA.ID, model.UpdateProductRequest{
		Price: &newPrice,
	})
	assert.ErrorIs(t, err, usecase.ErrProductAccessDenied)

	// 3. User B tries to delete User A's product -> MUST FAIL with ErrProductAccessDenied
	err = productUC.DeleteProduct(userB.User.ID, prodA.ID)
	assert.ErrorIs(t, err, usecase.ErrProductAccessDenied)

	// 4. User A legitimately updates own product -> MUST SUCCEED
	updatedPrice := float64(24000000)
	updatedProd, err := productUC.UpdateProduct(userA.User.ID, prodA.ID, model.UpdateProductRequest{
		Price: &updatedPrice,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedPrice, updatedProd.Price)

	// 5. Query with search and filter
	filterResult, total, err := productUC.GetAllProducts(model.ProductFilter{
		Search:     "Gaming",
		CategoryID: cat.ID,
		Limit:      10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.NotEmpty(t, filterResult)

	// 6. User A deletes own product -> MUST SUCCEED
	err = productUC.DeleteProduct(userA.User.ID, prodA.ID)
	require.NoError(t, err)

	// 7. Verify soft delete (product not found)
	_, err = productUC.GetProductByID(prodA.ID)
	assert.ErrorIs(t, err, usecase.ErrProductNotFound)
}
