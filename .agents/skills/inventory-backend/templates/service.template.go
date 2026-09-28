// internal/sale/service.go
//
// Services live INSIDE the owning domain package and own business logic plus
// transaction boundaries. They depend on repository INTERFACES — the domain's
// own and any cross-domain dependency — never on concrete GORM types or on
// other domains' services. They never import "gin" or return http.Status*.

package sale

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"i_m_s/internal/customer"
	"i_m_s/internal/inventory"
	"i_m_s/internal/models"
	"i_m_s/internal/product"
	"i_m_s/internal/utils/logger"
)

// Service depends on this domain's own Repository plus the interfaces of the
// domains it reads from. The dependency direction is one-way (sale -> inventory,
// product, customer); none of those may import sale. See design.md §3.1.
type Service struct {
	db        *gorm.DB
	sales     Repository
	products  product.Repository
	inventory inventory.Repository
	customers customer.Repository
}

func NewService(
	db *gorm.DB,
	sales Repository,
	products product.Repository,
	inventory inventory.Repository,
	customers customer.Repository,
) *Service {
	return &Service{db: db, sales: sales, products: products, inventory: inventory, customers: customers}
}

// CreateSale is the canonical example of AGENTS.md's transaction rule:
// sale + sale items + inventory decrement succeed or fail together.
func (s *Service) CreateSale(ctx context.Context, input CreateSaleInput) (*models.Sale, error) {
	contextLog := log.With().Str("customer_id", input.CustomerID.String()).Logger()
	contextLog.Info().Msg("creating sale")

	var sale *models.Sale

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Business rule: verify stock BEFORE writing anything.
		for _, item := range input.Items {
			inv, err := s.inventory.GetByProductID(ctx, tx, item.ProductID)
			if err != nil {
				return err // e.g. models.ErrNotFound, propagated as-is
			}
			if inv.Quantity < item.Quantity {
				contextLog.Warn().
					Str("product_id", item.ProductID.String()).
					Int("requested", item.Quantity).
					Int("available", inv.Quantity).
					Msg("sale rejected: insufficient inventory")
				return models.ErrInsufficientInventory
			}
		}

		// 2. Recompute totals server-side from trusted product data — never
		// trust client-supplied prices or totals.
		built, total, err := buildSaleFromInput(input)
		if err != nil {
			return err
		}
		built.TotalAmount = total

		if err := s.sales.Create(ctx, tx, built); err != nil {
			return err
		}

		// 3. Decrease inventory for each line item, same transaction.
		for _, item := range input.Items {
			if err := s.inventory.Decrease(ctx, tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}

		sale = built
		return nil
	})

	if err != nil {
		// logger.LogError (the package, not the local var above): full error
		// detail only in development; production logs a safe summary and marks
		// the detail suppressed. A silent transaction failure is a bug.
		logger.LogError(err, "sale transaction failed, rolled back", logger.Fields{
			"customer_id": input.CustomerID.String(),
		})
		return nil, err
	}

	contextLog.Info().
		Str("sale_id", sale.ID.String()).
		Int64("total_amount", sale.TotalAmount).
		Msg("sale completed")

	return sale, nil
}

type CreateSaleInput struct {
	CustomerID uuid.UUID
	UserID     uuid.UUID
	Items      []SaleItemInput
}

type SaleItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

// buildSaleFromInput constructs the domain object and computes totals from
// product records fetched via s.products. Omitted here for brevity — the point
// is that totals are computed in the service, never taken from the request.
func buildSaleFromInput(input CreateSaleInput) (*models.Sale, int64, error) {
	panic("implement: build *models.Sale + items and sum quantity*unit_price from product records")
}
