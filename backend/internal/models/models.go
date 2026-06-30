package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Roles string

const (
	Admin   Roles = "admin"
	Manager Roles = "manager"
	Staff   Roles = "staff"
)

type User struct {
	BaseModel
	
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	PassHash  string    `gorm:"not null" json:"-"`
	Role      Roles     `gorm:"not null" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}

type Supplier struct {
	BaseModel
	
	CompanyName   string    `gorm:"not null" json:"company_name"`
	ContactPerson string    `gorm:"not null" json:"contact_person"`
	Phone         string    `gorm:"not null" json:"phone"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Address       string    `gorm:"not null" json:"address"`
}

type Customer struct {
	BaseModel
	
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Phone     string    `gorm:"not null" json:"phone"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Address   string    `gorm:"not null" json:"address"`
}

type Category struct {
	BaseModel
	
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	BaseModel
	
	CategoryID  uuid.UUID `gorm:"type:uuid; not null" json:"category_id"`
	SKU         string    `gorm:"unique;not null" json:"sku"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	UnitPrice   int64     `gorm:"not null" json:"unit_price"`
	CostPrice   int64     `gorm:"not null" json:"cost_price"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

type Inventory struct {
	BaseModel
	
	ProductID    uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	ReorderLevel int       `gorm:"not null" json:"reorder_level"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

type PurchaseStatus string

const (
	PurchasePending   PurchaseStatus = "pending"
	PurchaseReceived  PurchaseStatus = "received"
	PurchaseCancelled PurchaseStatus = "cancelled"
)

type Purchase struct {
	BaseModel
	
	SupplierID   uuid.UUID      `gorm:"type:uuid; not null" json:"supplier_id"`
	UserID       uuid.UUID      `gorm:"type:uuid; not null" json:"user_id"`
	PurchaseDate time.Time      `gorm:"not null" json:"purchase_date"`
	Status       PurchaseStatus `gorm:"not null" json:"status"`
	TotalAmount  int64          `gorm:"not null" json:"total_amount"`
	Notes        string         `json:"notes"`

	Supplier Supplier       `gorm:"foreignKey:SupplierID" json:"supplier"`
	User     User           `gorm:"foreignKey:UserID" json:"-"`
	Items    []PurchaseItem `gorm:"foreignKey:PurchaseID"`
}

type PurchaseItem struct {
	BaseModel
	
	PurchaseID uuid.UUID `gorm:"type:uuid; not null" json:"purchase_id"`
	ProductID  uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	UnitCost   int64     `gorm:"not null" json:"unit_cost"`
	SubTotal   int64     `gorm:"not null" json:"subtotal"`

	Purchase Purchase `gorm:"foreignKey:PurchaseID" json:"purchase"`
	Product  Product  `gorm:"foreignKey:ProductID" json:"product"`
}

type SaleStatus string

const (
	SalePending   SaleStatus = "pending"
	SaleCompleted SaleStatus = "completed"
	SaleCancelled SaleStatus = "cancelled"
	SaleRefunded  SaleStatus = "refunded"
)

type PaymentMethod string

const (
	Cash         PaymentMethod = "cash"
	Card         PaymentMethod = "card"
	MobileMoney  PaymentMethod = "mobile_money"
	BankTransfer PaymentMethod = "bank_transfer"
)

type Sale struct {
	BaseModel
	
	CustomerID    uuid.UUID     `gorm:"type:uuid; not null" json:"customer_id"`
	UserID        uuid.UUID     `gorm:"type:uuid; not null" json:"user_id"`
	SaleDate      time.Time     `gorm:"not null" json:"sale_date"`
	Status        SaleStatus    `gorm:"not null" json:"status"`
	TotalAmount   int64         `gorm:"not null" json:"total_amount"`
	PaymentMethod PaymentMethod `gorm:"not null" json:"payment_method"`

	Customer Customer   `gorm:"foreignKey:CustomerID" json:"customer"`
	User     User       `gorm:"foreignKey:UserID" json:"-"`
	Items    []SaleItem `gorm:"foreignKey:SaleID"`
}

type SaleItem struct {
	BaseModel
	
	SaleID    uuid.UUID `gorm:"type:uuid; not null" json:"sale_id"`
	ProductID uuid.UUID `gorm:"type:uuid; not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice int64     `gorm:"not null" json:"unit_price"`
	SubTotal  int64     `gorm:"not null" json:"sub_total"`

	Sale    Sale    `gorm:"foreignKey:SaleID" json:"sale"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
