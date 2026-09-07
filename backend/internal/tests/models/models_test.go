package models_test

import (
	"testing"

	"i_m_s/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseModelBeforeCreateGeneratesUUID(t *testing.T) {
	b := models.BaseModel{}

	err := b.BeforeCreate(nil)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, b.ID)
}

func TestBaseModelPreservesExistingUUID(t *testing.T) {
	existingID := uuid.New()
	b := models.BaseModel{ID: existingID}

	err := b.BeforeCreate(nil)
	require.NoError(t, err)
	assert.Equal(t, existingID, b.ID)
}

func TestSentinelErrors(t *testing.T) {
	assert.Equal(t, "resource not found", models.ErrNotFound.Error())
	assert.Equal(t, "unauthorized", models.ErrUnauthorized.Error())
	assert.Equal(t, "forbidden", models.ErrForbidden.Error())
	assert.Equal(t, "insufficient inventory", models.ErrInsufficientInventory.Error())
	assert.Equal(t, "bad request", models.ErrBadRequest.Error())
	assert.Equal(t, "resource conflict", models.ErrConflict.Error())
	assert.Equal(t, "internal server error", models.ErrInternal.Error())
	assert.Equal(t, "validation failed", models.ErrValidationError.Error())
}