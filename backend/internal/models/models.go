package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Username  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"not null" json:"email"`
	PassHash  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null" json:"role"` // Admin, Manager, Staff
	IsActive  string    `gorm:"not null" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Supplier struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyName   string    `gorm:"not null" json:"company_name"`
	ContactPerson string    `gorm:"not null" json:"contact_person"`
	Phone         string    `gorm:"not null" json:"phone"`
	Email         string    `gorm:"not null" json:"email"`
	Address       string    `gorm:"not null" json:"address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_person"`
	Phone     string    `gorm:"not null" json:"phone"`
	Email     string    `gorm:"not null" json:"email"`
	Address   string    `gorm:"not null" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
}

type Product struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CategoryID  *uuid.UUID `gorm:"type:uuid; not null" json:"category_id"`
	SKU         string     `gorm:"not null" json:"sku"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"not null" json:"description"`
	UnitPrice   float64    `gorm:"not null" json:"unit_price"`
	CostPrice   float64    `gorm:"not null" json:"cost_price"`
	IsActive    string     `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Inventory struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    *uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity     int        `gorm:"not null" json:"quantity"`
	ReorderLevel int        `gorm:"not null" json:"reorder_level"`
	LastUpdated  time.Time  `gorm:"not null" json:"last_updated"`
}

type Purchase struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SupplierID   *uuid.UUID `gorm:"type:uuid; not null" json:"supplier_id"`
	UserID       *uuid.UUID `gorm:"type:uuid; not null" json:"user_id"`
	PurchaseDate time.Time  `gorm:"not null" json:"purchase_id"`
	Status       string     `gorm:"not null" json:"status"`
	TotalAmount  float64    `gorm:"not null" json:"total_amount"`
	Notes        string     `gorm:"null" json:"notes"`
}

type PurchaseItem struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PurchaseID *uuid.UUID `gorm:"type:uuid; not null" json:"purchase_id"`
	ProductID  *uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	UnitCost   float64    `gorm:"not null" json:"unit_cost"`
	SubTotal   float64    `gorm:"not null" json:"subtotal"`
}

type Sales struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CustomerID    *uuid.UUID `gorm:"type:uuid; not null" json:"customer_id"`
	UserID        *uuid.UUID `gorm:"type:uuid; not null" json:"user_id"`
	SaleDate      time.Time  `gorm:"not null" json:"sale_date"`
	Status        string     `gorm:"not null" json:"status"`
	TotalAmount   float64    `gorm:"not null" json:"total_amount"`
	PaymentMethod string     `gorm:"not null" json:"payment_method"`
}

type SaleItem struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SaleID    *uuid.UUID `gorm:"type:uuid; not null" json:"sale_id"`
	ProductID *uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity  int        `gorm:"not null" json:"quantity"`
	UnitPrice float64    `gorm:"not null" json:"unit_price"`
	SubTotal  float64    `gorm:"not null" json:"sub_total"`
}
