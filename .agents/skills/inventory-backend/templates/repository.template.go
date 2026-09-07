// internal/repositories/product_repository.go
//
// Repositories: pure persistence. No business logic, no HTTP, no JWT checks.
// Accept *gorm.DB (not a wrapped struct with its own connection) so the same
// repository works both standalone and inside a service's transaction.

package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"inventory/internal/models"
)

// Interface, not just a struct — this is what lets services depend on an
// abstraction (DIP) and lets tests substitute a mock without hitting Postgres.
type ProductRepository interface {
	Create(ctx context.Context, db *gorm.DB, product *models.Product) error
	GetByID(ctx context.Context, db *gorm.DB, id string) (*models.Product, error)
	Update(ctx context.Context, db *gorm.DB, product *models.Product) error
	ListByCategory(ctx context.Context, db *gorm.DB, categoryID string) ([]models.Product, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) Create(ctx context.Context, db *gorm.DB, product *models.Product) error {
	if err := db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, db *gorm.DB, id string) (*models.Product, error) {
	var product models.Product
	err := db.WithContext(ctx).First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Translate GORM's error to our own sentinel so services never
		// need to import gorm's error types.
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product %s: %w", id, err)
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, db *gorm.DB, product *models.Product) error {
	if err := db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("update product %s: %w", product.ID, err)
	}
	return nil
}

func (r *productRepository) ListByCategory(ctx context.Context, db *gorm.DB, categoryID string) ([]models.Product, error) {
	var products []models.Product
	// Load only what the use case needs — no blanket Preload here.
	err := db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("list products for category %s: %w", categoryID, err)
	}
	return products, nil
}
