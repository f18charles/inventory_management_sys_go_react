// internal/services/sale_service_test.go
//
// Service tests mock the repository INTERFACE — this is the concrete payoff
// of depending on interfaces (DIP) rather than concrete GORM types. Assert
// business outcomes, not that method X called method Y.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"inventory/internal/models"
	"inventory/internal/services"
)

// Mock satisfying the InventoryRepository interface.
type mockInventoryRepo struct{ mock.Mock }

func (m *mockInventoryRepo) GetByProductID(ctx context.Context, db *gorm.DB, productID string) (*models.Inventory, error) {
	args := m.Called(ctx, db, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Inventory), args.Error(1)
}

func (m *mockInventoryRepo) Decrease(ctx context.Context, db *gorm.DB, productID string, qty int) error {
	args := m.Called(ctx, db, productID, qty)
	return args.Error(0)
}

func TestCreateSale_RejectsWhenInsufficientStock(t *testing.T) {
	// Arrange: stock is 2, sale requests 5.
	invRepo := new(mockInventoryRepo)
	invRepo.On("GetByProductID", mock.Anything, mock.Anything, "product-1").
		Return(&models.Inventory{ProductID: "product-1", Quantity: 2}, nil)

	svc := services.NewSaleService(nil, nil, nil, invRepo) // db/saleRepo/productRepo omitted for this table case

	// Act
	_, err := svc.CreateSale(context.Background(), services.CreateSaleInput{
		CustomerID: "customer-1",
		Items: []services.SaleItemInput{
			{ProductID: "product-1", Quantity: 5},
		},
	})

	// Assert: the business rule was enforced, not just "no panic".
	require.Error(t, err)
	assert.ErrorIs(t, err, models.ErrInsufficientInventory)
	invRepo.AssertNotCalled(t, "Decrease", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Table-driven pattern for multiple cases on the same behavior:
func TestCreateSale_StockCases(t *testing.T) {
	cases := []struct {
		name      string
		available int
		requested int
		wantErr   error
	}{
		{"exact match succeeds", 5, 5, nil},
		{"more than available fails", 3, 5, models.ErrInsufficientInventory},
		{"zero available fails", 0, 1, models.ErrInsufficientInventory},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// ...set up mocks per-case using tc.available / tc.requested,
			// then assert errors.Is(err, tc.wantErr) or require.NoError(t, err).
		})
	}
}
