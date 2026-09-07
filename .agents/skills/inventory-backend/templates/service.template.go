// internal/services/sale_service.go
//
// Services: business logic + transaction boundaries. Depend on repository
// INTERFACES (never a concrete GORM type) so tests can mock them. Never
// import "gin" or reference *gin.Context. Never return http.Status*.

package services

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"inventory/internal/models"
	"inventory/internal/repositories"
	"inventory/internal/utils/logger"
)

type SaleService struct {
	db            *gorm.DB
	saleRepo      repositories.SaleRepository
	productRepo   repositories.ProductRepository
	inventoryRepo repositories.InventoryRepository
}

// SaleRepository and InventoryRepository follow the exact same interface
// pattern as ProductRepository in templates/repository.template.go — an
// interface + a plain struct implementation taking (ctx, db, ...args).

func NewSaleService(
	db *gorm.DB,
	saleRepo repositories.SaleRepository,
	productRepo repositories.ProductRepository,
	inventoryRepo repositories.InventoryRepository,
) *SaleService {
	return &SaleService{db: db, saleRepo: saleRepo, productRepo: productRepo, inventoryRepo: inventoryRepo}
}

// CreateSale is the canonical example of AGENTS.md's transaction rule:
// sale + sale items + inventory decrement succeed or fail together.
func (s *SaleService) CreateSale(ctx context.Context, input CreateSaleInput) (*models.Sale, error) {
	contextLog := log.With().Str("customer_id", input.CustomerID).Logger()
	contextLog.Info().Msg("creating sale")

	var sale *models.Sale

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Business rule: verify stock BEFORE writing anything.
		for _, item := range input.Items {
			inv, err := s.inventoryRepo.GetByProductID(ctx, tx, item.ProductID)
			if err != nil {
				return err // e.g. models.ErrNotFound, propagated as-is
			}
			if inv.Quantity < item.Quantity {
				contextLog.Warn().
					Str("product_id", item.ProductID).
					Int("requested", item.Quantity).
					Int("available", inv.Quantity).
					Msg("sale rejected: insufficient inventory")
				return models.ErrInsufficientInventory
			}
		}

		// 2. Recompute totals server-side — never trust client-provided totals.
		built, total := buildSaleFromInput(input)
		built.TotalAmount = total

		if err := s.saleRepo.Create(ctx, tx, built); err != nil {
			return err
		}

		// 3. Decrease inventory for each line item, same transaction.
		for _, item := range input.Items {
			if err := s.inventoryRepo.Decrease(ctx, tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}

		sale = built
		return nil
	})

	if err != nil {
		// logger.LogError (the package, not the local var above): full err
		// detail + stack only in development; production logs a safe
		// summary and marks detail as suppressed. See internal/utils/logger
		// — a silent transaction failure is treated as a bug.
		logger.LogError(err, "sale transaction failed, rolled back", logger.Fields{
			"customer_id": input.CustomerID,
		})
		return nil, err
	}

	contextLog.Info().
		Str("sale_id", sale.ID).
		Int64("total_amount", sale.TotalAmount).
		Msg("sale completed")

	return sale, nil
}

type CreateSaleInput struct {
	CustomerID string
	UserID     string
	Items      []SaleItemInput
}

type SaleItemInput struct {
	ProductID string
	Quantity  int
}

// buildSaleFromInput and total calculation omitted for brevity in this
// template — the point being illustrated is: build the domain object and
// compute totals here, in the service, not from client-supplied numbers.
func buildSaleFromInput(input CreateSaleInput) (*models.Sale, int64) {
	panic("implement: construct *models.Sale + sale items, sum quantity*unit_price from product cost data")
}
