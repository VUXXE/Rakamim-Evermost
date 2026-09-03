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

func TestTransactionUseCase_AtomicCheckout_StockDeduction_AndSnapshots(t *testing.T) {
	ctx := setupAllTestDB(t)

	categoryRepo := repository.NewCategoryRepository(ctx.db)
	productRepo := repository.NewProductRepository(ctx.db)
	addressRepo := repository.NewAddressRepository(ctx.db)
	txRepo := repository.NewTransactionRepository(ctx.db)
	logRepo := repository.NewProductLogRepository(ctx.db)

	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	productUC := usecase.NewProductUseCase(productRepo, ctx.storeRepo, categoryRepo)
	addressUC := usecase.NewAddressUseCase(ctx.db, addressRepo)
	txUC := usecase.NewTransactionUseCase(ctx.db, txRepo, logRepo, productRepo, addressRepo)

	ts := time.Now().UnixNano()

	// 1. Setup Seller & Product
	seller, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Seller Merchant",
		Email:    fmt.Sprintf("merchant_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0871%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	cat, err := categoryUC.CreateCategory(model.CreateCategoryRequest{
		Name: fmt.Sprintf("Apparel_%d", ts%100000),
	})
	require.NoError(t, err)

	prod, err := productUC.CreateProduct(seller.User.ID, model.CreateProductRequest{
		CategoryID:  cat.ID,
		Name:        "Evermos Hoodie",
		Description: "Warm hoodie",
		Price:       250000,
		Quantity:    5, // Initial stock: 5
	})
	require.NoError(t, err)

	// 2. Setup Buyer & Address
	buyer, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Buyer Customer",
		Email:    fmt.Sprintf("buyer_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0872%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	addr, err := addressUC.CreateAddress(buyer.User.ID, model.CreateAddressRequest{
		JudulAlamat:   "Rumah Pembeli",
		PenerimaNama:  "Buyer Customer",
		PenerimaPhone: "0872000000",
		Provinsi:      "DKI Jakarta",
		ProvinsiID:    "31",
		Kabupaten:     "Jakarta Selatan",
		KabupatenID:   "3174",
		Kecamatan:     "Tebet",
		KecamatanID:   "317402",
		Kelurahan:     "Tebet Barat",
		KelurahanID:   "3174021001",
		DetailAlamat:  "Jl. Tebet Barat No. 5",
		IsDefault:     true,
	})
	require.NoError(t, err)

	// 3. Test Checkout with Insufficient Stock -> MUST FAIL and ROLLBACK
	_, err = txUC.CreateTransaction(buyer.User.ID, model.CreateTransactionRequest{
		AddressID: addr.ID,
		Products: []model.CheckoutItem{
			{ProductID: prod.ID, Quantity: 10}, // More than stock (5)
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")

	// Verify stock was not decremented
	refetchedProd, err := productUC.GetProductByID(prod.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, refetchedProd.Quantity)

	// 4. Test Successful Checkout (Purchase 2 hoodies)
	txResp, err := txUC.CreateTransaction(buyer.User.ID, model.CreateTransactionRequest{
		AddressID: addr.ID,
		Products: []model.CheckoutItem{
			{ProductID: prod.ID, Quantity: 2},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "pending", txResp.Status)
	assert.Equal(t, float64(500000), txResp.TotalPrice)
	assert.NotEmpty(t, txResp.InvoiceNumber)
	assert.Len(t, txResp.ProductLogs, 1)
	assert.Equal(t, prod.ID, txResp.ProductLogs[0].ProductID)
	assert.Equal(t, 2, txResp.ProductLogs[0].Quantity)
	assert.Equal(t, float64(250000), txResp.ProductLogs[0].Price)

	// Verify stock was decremented from 5 to 3
	refetchedProd2, err := productUC.GetProductByID(prod.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, refetchedProd2.Quantity)

	// 5. Test Buyer fetches own transaction
	fetchedTx, err := txUC.GetMyTransactionByID(buyer.User.ID, txResp.ID)
	require.NoError(t, err)
	assert.Equal(t, txResp.InvoiceNumber, fetchedTx.InvoiceNumber)
	assert.Len(t, fetchedTx.ProductLogs, 1)

	// 6. Test Multi-Tenant Isolation: Another user cannot access Buyer's transaction
	anotherUser, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Another User",
		Email:    fmt.Sprintf("intruder_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0873%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	_, err = txUC.GetMyTransactionByID(anotherUser.User.ID, txResp.ID)
	assert.ErrorIs(t, err, usecase.ErrTransactionNotFound)

	// 7. Buyer updates status: "completed"
	completedTx, err := txUC.UpdateTransactionStatus(buyer.User.ID, txResp.ID, "completed", false)
	require.NoError(t, err)
	assert.Equal(t, "completed", completedTx.Status)
}
