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
)

func setupTestDB(t *testing.T) usecase.AuthUseCase {
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
	require.NoError(t, err, "failed to connect to test MySQL")

	err = mysql.AutoMigrate(db)
	require.NoError(t, err, "auto-migrate failed")

	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)

	return usecase.NewAuthUseCase(db, userRepo, storeRepo)
}

func TestAuthUseCase_Register_Success(t *testing.T) {
	authUC := setupTestDB(t)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("test_%d@evermos.com", timestamp)
	phone := fmt.Sprintf("0812%d", timestamp%100000000)
	name := fmt.Sprintf("Merchant_%d", timestamp%10000)

	req := model.RegisterRequest{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: "SecurePassword123!",
	}

	resp, err := authUC.Register(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify User
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, email, resp.User.Email)
	assert.Equal(t, phone, resp.User.Phone)
	assert.Equal(t, name, resp.User.Name)
	assert.False(t, resp.User.IsAdmin)

	// Verify Invariant #1: Automatic Store Creation
	require.NotNil(t, resp.Store)
	assert.Equal(t, resp.User.ID, resp.Store.UserID)
	expectedStoreName := fmt.Sprintf("%s's Store", name)
	assert.Equal(t, expectedStoreName, resp.Store.Name)
}

func TestAuthUseCase_Register_DuplicateEmail(t *testing.T) {
	authUC := setupTestDB(t)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("dup_email_%d@evermos.com", timestamp)
	phone1 := fmt.Sprintf("0811%d", timestamp%100000000)
	phone2 := fmt.Sprintf("0812%d", timestamp%100000000)

	req1 := model.RegisterRequest{
		Name:     "User One",
		Email:    email,
		Phone:    phone1,
		Password: "Password123!",
	}

	_, err := authUC.Register(req1)
	require.NoError(t, err)

	req2 := model.RegisterRequest{
		Name:     "User Two",
		Email:    email,
		Phone:    phone2,
		Password: "Password123!",
	}

	_, err = authUC.Register(req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already registered")
}

func TestAuthUseCase_Register_DuplicatePhone(t *testing.T) {
	authUC := setupTestDB(t)

	timestamp := time.Now().UnixNano()
	phone := fmt.Sprintf("0813%d", timestamp%100000000)
	email1 := fmt.Sprintf("user1_%d@evermos.com", timestamp)
	email2 := fmt.Sprintf("user2_%d@evermos.com", timestamp)

	req1 := model.RegisterRequest{
		Name:     "User One",
		Email:    email1,
		Phone:    phone,
		Password: "Password123!",
	}

	_, err := authUC.Register(req1)
	require.NoError(t, err)

	req2 := model.RegisterRequest{
		Name:     "User Two",
		Email:    email2,
		Phone:    phone,
		Password: "Password123!",
	}

	_, err = authUC.Register(req2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone already registered")
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	authUC := setupTestDB(t)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("login_success_%d@evermos.com", timestamp)
	phone := fmt.Sprintf("0814%d", timestamp%100000000)
	password := "CorrectPassword123!"

	_, err := authUC.Register(model.RegisterRequest{
		Name:     "Login User",
		Email:    email,
		Phone:    phone,
		Password: password,
	})
	require.NoError(t, err)

	loginResp, err := authUC.Login(model.LoginRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	require.NotNil(t, loginResp)
	assert.NotEmpty(t, loginResp.Token)
	assert.Equal(t, email, loginResp.User.Email)
	assert.NotNil(t, loginResp.Store)
	assert.Equal(t, "Login User's Store", loginResp.Store.Name)
}

func TestAuthUseCase_Login_InvalidPassword(t *testing.T) {
	authUC := setupTestDB(t)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("login_fail_%d@evermos.com", timestamp)
	phone := fmt.Sprintf("0815%d", timestamp%100000000)

	_, err := authUC.Register(model.RegisterRequest{
		Name:     "Wrong Password User",
		Email:    email,
		Phone:    phone,
		Password: "ActualPassword123!",
	})
	require.NoError(t, err)

	_, err = authUC.Login(model.LoginRequest{
		Email:    email,
		Password: "WrongPassword!",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}

func TestAuthUseCase_Login_NonExistentEmail(t *testing.T) {
	authUC := setupTestDB(t)

	_, err := authUC.Login(model.LoginRequest{
		Email:    "nonexistent@evermos.com",
		Password: "AnyPassword123!",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}
