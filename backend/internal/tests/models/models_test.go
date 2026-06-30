package models_test

import (
	"i_m_s/internal/models"
	"testing"

	"github.com/google/uuid"
)

func TestBaseModelBeforeCreateGeneratesUUID(t *testing.T) {
	b := models.BaseModel{}

	err := b.BeforeCreate(nil)

	if err != nil {
		t.Fatal(err)
	}

	if b.ID == uuid.Nil {
		t.Fatal("expected UUID to be generated")
	}
}