package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"evermos-api/internal/infrastructure/mysql"
	"evermos-api/internal/pkg/model"
	"evermos-api/internal/pkg/repository"
	"evermos-api/internal/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type testContext struct {
	db        *gorm.DB
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
	authUC    usecase.AuthUseCase
	userUC    usecase.UserUseCase
	storeUC   usecase.StoreUseCase
}

func setupAllTestDB(t *testing.T) *testContext {
	cfg := mysql.Config{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Password:        "password",
		DBName:          "evermos",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := mysql.Connect(cfg)
	require.NoError(t, err)

	err = mysql.AutoMigrate(db)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)

	authUC := usecase.NewAuthUseCase(db, userRepo, storeRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	storeUC := usecase.NewStoreUseCase(storeRepo)

	return &testContext{
		db:        db,
		userRepo:  userRepo,
		storeRepo: storeRepo,
		authUC:    authUC,
		userUC:    userUC,
		storeUC:   storeUC,
	}
}

func TestUserUseCase_GetMe_AndUpdateMe(t *testing.T) {
	ctx := setupAllTestDB(t)

	ts := time.Now().UnixNano()
	authResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Profile User",
		Email:    fmt.Sprintf("profile_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0877%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// GetMe
	user, err := ctx.userUC.GetMe(authResp.User.ID)
	require.NoError(t, err)
	assert.Equal(t, "Profile User", user.Name)

	// UpdateMe
	newName := "Profile User Updated"
	updated, err := ctx.userUC.UpdateMe(authResp.User.ID, model.UpdateUserRequest{
		Name: newName,
	})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
}

func TestStoreUseCase_GetMyStore(t *testing.T) {
	ctx := setupAllTestDB(t)

	ts := time.Now().UnixNano()
	name := fmt.Sprintf("Merchant_%d", ts%10000)
	authResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     name,
		Email:    fmt.Sprintf("merchant_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0878%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	store, err := ctx.storeUC.GetMyStore(authResp.User.ID)
	require.NoError(t, err)
	require.NotNil(t, store)
	assert.Equal(t, fmt.Sprintf("%s's Store", name), store.Name)
}

func TestStoreUseCase_MultiTenantIsolation_UpdateStore(t *testing.T) {
	ctx := setupAllTestDB(t)

	ts := time.Now().UnixNano()
	// User A
	userAResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "User A",
		Email:    fmt.Sprintf("user_a_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0871%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// User B
	userBResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "User B",
		Email:    fmt.Sprintf("user_b_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0872%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// User A attempts to update User B's store -> MUST FAIL with ErrStoreAccessDenied
	_, err = ctx.storeUC.UpdateStore(userAResp.User.ID, userBResp.Store.ID, model.UpdateStoreRequest{
		Name: "Hacked Store Name",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrStoreAccessDenied)

	// User B legitimately updates own store -> MUST SUCCEED
	updatedStore, err := ctx.storeUC.UpdateStore(userBResp.User.ID, userBResp.Store.ID, model.UpdateStoreRequest{
		Name:    "User B's Official Store",
		Address: "Jakarta Pusat",
	})
	require.NoError(t, err)
	assert.Equal(t, "User B's Official Store", updatedStore.Name)
	assert.Equal(t, "Jakarta Pusat", updatedStore.Address)
}

func TestStoreUseCase_MultiTenantIsolation_DeleteStore(t *testing.T) {
	ctx := setupAllTestDB(t)

	ts := time.Now().UnixNano()
	userAResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Owner A",
		Email:    fmt.Sprintf("owner_a_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0873%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	userBResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Owner B",
		Email:    fmt.Sprintf("owner_b_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0874%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// User A attempts to delete User B's store -> MUST FAIL
	err = ctx.storeUC.DeleteStore(userAResp.User.ID, userBResp.Store.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrStoreAccessDenied)

	// User B deletes own store -> MUST SUCCEED
	err = ctx.storeUC.DeleteStore(userBResp.User.ID, userBResp.Store.ID)
	assert.NoError(t, err)
}
