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

func TestAddressUseCase_Create_DefaultAddressHandling(t *testing.T) {
	ctx := setupAllTestDB(t)
	addressRepo := repository.NewAddressRepository(ctx.db)
	addressUC := usecase.NewAddressUseCase(ctx.db, addressRepo)

	ts := time.Now().UnixNano()
	authResp, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "Address User",
		Email:    fmt.Sprintf("address_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0891%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	// Address 1: Default
	addr1, err := addressUC.CreateAddress(authResp.User.ID, model.CreateAddressRequest{
		JudulAlamat:   "Rumah",
		PenerimaNama:  "John Doe",
		PenerimaPhone: "081234567890",
		Provinsi:      "DKI Jakarta",
		ProvinsiID:    "31",
		Kabupaten:     "Jakarta Selatan",
		KabupatenID:   "3174",
		Kecamatan:     "Kebayoran Baru",
		KecamatanID:   "317401",
		Kelurahan:     "Selong",
		KelurahanID:   "3174011001",
		DetailAlamat:  "Jl. Senopati No. 10",
		IsDefault:     true,
	})
	require.NoError(t, err)
	assert.True(t, addr1.IsDefault)

	// Address 2: Also created as Default -> MUST unset Address 1
	addr2, err := addressUC.CreateAddress(authResp.User.ID, model.CreateAddressRequest{
		JudulAlamat:   "Kantor",
		PenerimaNama:  "John Doe",
		PenerimaPhone: "081234567890",
		Provinsi:      "DKI Jakarta",
		ProvinsiID:    "31",
		Kabupaten:     "Jakarta Pusat",
		KabupatenID:   "3171",
		Kecamatan:     "Gambir",
		KecamatanID:   "317101",
		Kelurahan:     "Gambir",
		KelurahanID:   "3171011001",
		DetailAlamat:  "Jl. Medan Merdeka Barat No. 1",
		IsDefault:     true,
	})
	require.NoError(t, err)
	assert.True(t, addr2.IsDefault)

	// Re-fetch Address 1: MUST now be false
	refetchedAddr1, err := addressUC.GetAddressByID(authResp.User.ID, addr1.ID)
	require.NoError(t, err)
	assert.False(t, refetchedAddr1.IsDefault)
}

func TestAddressUseCase_MultiTenantIsolation(t *testing.T) {
	ctx := setupAllTestDB(t)
	addressRepo := repository.NewAddressRepository(ctx.db)
	addressUC := usecase.NewAddressUseCase(ctx.db, addressRepo)

	ts := time.Now().UnixNano()
	userA, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "User A",
		Email:    fmt.Sprintf("user_a_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0892%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	userB, err := ctx.authUC.Register(model.RegisterRequest{
		Name:     "User B",
		Email:    fmt.Sprintf("user_b_%d@evermos.com", ts),
		Phone:    fmt.Sprintf("0893%d", ts%100000000),
		Password: "Password123!",
	})
	require.NoError(t, err)

	addrA, err := addressUC.CreateAddress(userA.User.ID, model.CreateAddressRequest{
		JudulAlamat:   "Alamat User A",
		PenerimaNama:  "User A",
		PenerimaPhone: "0892000000",
		Provinsi:      "Jawa Barat",
		ProvinsiID:    "32",
		Kabupaten:     "Bandung",
		KabupatenID:   "3273",
		Kecamatan:     "Coblong",
		KecamatanID:   "327301",
		Kelurahan:     "Dago",
		KelurahanID:   "3273011001",
		DetailAlamat:  "Jl. Dago No. 1",
		IsDefault:     false,
	})
	require.NoError(t, err)

	// User B attempts to access User A's address -> MUST FAIL with ErrAddressNotFound
	_, err = addressUC.GetAddressByID(userB.User.ID, addrA.ID)
	assert.ErrorIs(t, err, usecase.ErrAddressNotFound)

	// User B attempts to update User A's address -> MUST FAIL with ErrAddressNotFound
	newTitle := "Hacked Address Title"
	_, err = addressUC.UpdateAddress(userB.User.ID, addrA.ID, model.UpdateAddressRequest{
		JudulAlamat: newTitle,
	})
	assert.ErrorIs(t, err, usecase.ErrAddressNotFound)

	// User B attempts to delete User A's address -> MUST FAIL with ErrAddressNotFound
	err = addressUC.DeleteAddress(userB.User.ID, addrA.ID)
	assert.ErrorIs(t, err, usecase.ErrAddressNotFound)

	// User A legitimately updates own address -> MUST SUCCEED
	updated, err := addressUC.UpdateAddress(userA.User.ID, addrA.ID, model.UpdateAddressRequest{
		JudulAlamat: "Alamat Baru A",
	})
	require.NoError(t, err)
	assert.Equal(t, "Alamat Baru A", updated.JudulAlamat)
}
