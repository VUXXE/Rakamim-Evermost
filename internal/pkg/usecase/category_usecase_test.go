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

func TestCategoryUseCase_CRUD(t *testing.T) {
	ctx := setupAllTestDB(t)
	categoryRepo := repository.NewCategoryRepository(ctx.db)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)

	ts := time.Now().UnixNano()
	catName := fmt.Sprintf("Electronics_%d", ts%100000)

	// 1. Create Category
	created, err := categoryUC.CreateCategory(model.CreateCategoryRequest{
		Name: catName,
	})
	require.NoError(t, err)
	assert.Equal(t, catName, created.Name)

	// 2. Duplicate Name -> Error
	_, err = categoryUC.CreateCategory(model.CreateCategoryRequest{
		Name: catName,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// 3. Get By ID
	fetched, err := categoryUC.GetCategoryByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, catName, fetched.Name)

	// 4. Update
	updatedName := catName + "_Updated"
	updated, err := categoryUC.UpdateCategory(created.ID, model.UpdateCategoryRequest{
		Name: updatedName,
	})
	require.NoError(t, err)
	assert.Equal(t, updatedName, updated.Name)

	// 5. Delete
	err = categoryUC.DeleteCategory(created.ID)
	require.NoError(t, err)

	// 6. Verify NotFound after delete
	_, err = categoryUC.GetCategoryByID(created.ID)
	assert.ErrorIs(t, err, usecase.ErrCategoryNotFound)
}
