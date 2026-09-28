// internal/product/repository.go
//
// Repositories live INSIDE the owning domain package (design.md §3.2) and are
// pure persistence: CRUD + queries only. No business logic, no HTTP, no JWT.
// Each method accepts *gorm.DB so the same repository works standalone or
// inside a transaction owned by the service.

package product

import (
	"context"
	"errors"
	"fmt"

	"i_m_s/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository is the domain's persistence boundary. Other domains and this
// domain's service depend on this interface (DIP), which is what lets service
// tests substitute a mock without hitting PostgreSQL.
type Repository interface {
	Create(ctx context.Context, db *gorm.DB, product *models.Product) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.Product, error)
	Update(ctx context.Context, db *gorm.DB, product *models.Product) error
	ListByCategory(ctx context.Context, db *gorm.DB, categoryID uuid.UUID) ([]models.Product, error)
}

type repository struct{}

// NewRepository returns the GORM-backed implementation of Repository.
func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, db *gorm.DB, product *models.Product) error {
	err := db.WithContext(ctx).Create(product).Error
	// TranslateError is enabled on the connection, so unique violations arrive
	// as gorm.ErrDuplicatedKey; map them to a domain sentinel and never let a
	// driver error escape the repository.
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return models.ErrConflict
	}
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := db.WithContext(ctx).First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Translate GORM's error to our sentinel so services never import
		// gorm's error types.
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product %s: %w", id, err)
	}
	return &product, nil
}

func (r *repository) Update(ctx context.Context, db *gorm.DB, product *models.Product) error {
	if err := db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("update product %s: %w", product.ID, err)
	}
	return nil
}

func (r *repository) ListByCategory(ctx context.Context, db *gorm.DB, categoryID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	// Load only what the use case needs — no blanket Preload here.
	if err := db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("list products for category %s: %w", categoryID, err)
	}
	return products, nil
}
