# Design Document: Inventory Management System

This document records the architecture, schema, and technical decisions for the Inventory Management System. It ensures consistent engineering choices across the backend and frontend. Read this together with `PRD.md` (product requirements and scope) and `AGENTS.md` (development ground rules).

---

## 1. Technology Stack

| Layer | Choice | Rationale |
|---|---|---|
| **Backend Framework** | Go with Gin (`gin-gonic/gin`) | High throughput, low memory footprint, simple middleware pipeline, and clear handler lifecycle. |
| **ORM & Data Mapping** | GORM (`gorm.io/gorm`) | Idiomatic Go struct mapping, query building, relation preloading, and migration compatibility. |
| **Database** | PostgreSQL 16+ | ACID compliance, native UUIDs, foreign key constraints, check constraints, and row-level locking. |
| **Database Migrations** | `golang-migrate/migrate` | Versioned SQL files in `backend/internal/database/migrations` (`000001_*.up.sql` / `.down.sql`). |
| **Authentication & Crypto** | JWT (`golang-jwt/jwt/v5`) & bcrypt (`golang.org/x/crypto/bcrypt`) | Stateless authorization tokens with standard claims; secure, salt-backed password hashing. |
| **Logging** | `zerolog` (`rs/zerolog`) | High-performance structured JSON logging with context fields and environment-driven formatting. |
| **API Documentation** | Swagger / OpenAPI (`swaggo/swag`, `gin-swagger`) | Code-generated API specifications directly from Gin handler annotations. |
| **Testing** | Go standard `testing` + Testify (`stretchr/testify`) | Assertion suites, mocks, and integration test helpers. |
| **Frontend Framework** | React 18+ (JavaScript / JSX) with Vite | Fast development builds, optimized bundle size, and mature component ecosystem. |
| **Frontend Routing** | TanStack Router | Type-safe client-side routing, nested layouts, and route loaders. |
| **Server State & Caching** | TanStack Query (`@tanstack/react-query`) | Declarative API fetching, background cache revalidation, mutation state, and optimistic updates. |
| **Client State** | Zustand | Minimalist, boilerplate-free state store for session auth tokens and ephemeral UI state. |
| **HTTP Client** | Axios | Interceptors for JWT attachment, error transformation, and standard envelope unpacking. |
| **Styling** | Tailwind CSS | Utility-first responsive design, consistent design tokens, and fast build times. |

---

## 2. Data Model & Schema (v1)

### 2.1. Base Entity & Conventions
Every entity inherits `BaseModel` (defined in `backend/internal/models/models.go`):
- `id`: UUID primary key, generated automatically via GORM `BeforeCreate` hook if `uuid.Nil`.
- `created_at`: `TIMESTAMPTZ`, default `CURRENT_TIMESTAMP`.
- `updated_at`: `TIMESTAMPTZ`, default `CURRENT_TIMESTAMP`.
- `deleted_at`: `TIMESTAMPTZ`, soft deletion managed by GORM.

Monetary values (`unit_price`, `cost_price`, `unit_cost`, `subtotal`, `total_amount`) are stored as `int64` (cents/minor units). Stock quantities are non-negative integers (`int`).

### 2.2. Entity Relationships
```
Category (1) ───< (N) Product (1) ─── (1) Inventory
Supplier (1) ───< (N) Purchase (1) ───< (N) PurchaseItem (N) >─── (1) Product
Customer (1) ───< (N) Sale (1) ───────< (N) SaleItem (N) >─────── (1) Product
User (1) ───────< (N) Purchases, Sales
```

### 2.3. Schema Definition (PostgreSQL DDL matching `backend/internal/database/migrations`)

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    pass_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'staff')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_users_username_unique ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email_unique ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Suppliers table
CREATE TABLE suppliers (
    id UUID PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_person TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_suppliers_email_unique ON suppliers(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_deleted_at ON suppliers(deleted_at);

-- Customers table
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_customers_email_unique ON customers(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

-- Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_categories_name_unique ON categories(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    sku TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    unit_price BIGINT NOT NULL CHECK (unit_price >= 0),
    cost_price BIGINT NOT NULL CHECK (cost_price >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_products_sku_unique ON products(sku) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

-- Inventories table
CREATE TABLE inventories (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reorder_level INT NOT NULL DEFAULT 0 CHECK (reorder_level >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_inventories_product_id_unique ON inventories(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_inventories_deleted_at ON inventories(deleted_at);

-- Purchases table
CREATE TABLE purchases (
    id UUID PRIMARY KEY,
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    purchase_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('pending', 'received', 'cancelled')),
    total_amount BIGINT NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_purchases_supplier_id ON purchases(supplier_id);
CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_purchases_status ON purchases(status);
CREATE INDEX idx_purchases_deleted_at ON purchases(deleted_at);

-- Purchase Items table
CREATE TABLE purchase_items (
    id UUID PRIMARY KEY,
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_cost BIGINT NOT NULL CHECK (unit_cost >= 0),
    subtotal BIGINT NOT NULL CHECK (subtotal >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX idx_purchase_items_product_id ON purchase_items(product_id);
CREATE INDEX idx_purchase_items_deleted_at ON purchase_items(deleted_at);

-- Sales table
CREATE TABLE sales (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    sale_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'cancelled', 'refunded')),
    total_amount BIGINT NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'card', 'mobile_money', 'bank_transfer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_sales_customer_id ON sales(customer_id);
CREATE INDEX idx_sales_user_id ON sales(user_id);
CREATE INDEX idx_sales_status ON sales(status);
CREATE INDEX idx_sales_deleted_at ON sales(deleted_at);

-- Sale Items table
CREATE TABLE sale_items (
    id UUID PRIMARY KEY,
    sale_id UUID NOT NULL REFERENCES sales(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price BIGINT NOT NULL CHECK (unit_price >= 0),
    sub_total BIGINT NOT NULL CHECK (sub_total >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_sale_items_product_id ON sale_items(product_id);
CREATE INDEX idx_sale_items_deleted_at ON sale_items(deleted_at);
```

---

## 3. Package Architecture: Domain-Based (Package by Feature)

To keep code cohesive, scalable, and modular, the codebase organizes logic into domain packages under `internal/`.

### 3.1. Preventing Cyclic Import Errors in Go
In Go, package cycles (`package A imports B` and `package B imports A`) cause compile-time failures. To eliminate cyclic dependencies completely:

1. **Shared Domain Entities & Enums in `internal/models`**:
   - All GORM database models (`User`, `Category`, `Supplier`, `Customer`, `Product`, `Inventory`, `Purchase`, `PurchaseItem`, `Sale`, `SaleItem`) and shared types (`Roles`, `PurchaseStatus`, `SaleStatus`, `PaymentMethod`) live in `backend/internal/models`.
   - `internal/models` imports only the standard library, `google/uuid`, and `gorm.io/gorm`. It never imports any domain package.
   - All domain packages import `internal/models` without creating cycles.

2. **Acyclic Domain Dependencies (Directed Acyclic Graph)**:
   - When a domain service requires functionality from another domain (e.g., `sale.Service` checking stock or updating inventory), it depends on the target domain's interface (e.g. `inventory.Service` or `inventory.Repository`).
   - The dependency flow is strictly one-directional:
     ```
     sale       ──> inventory, product, customer
     purchase   ──> inventory, product, supplier
     inventory  ──> product
     product    ──> category
     auth       ──> user
     report     ──> sale, purchase, inventory, product
     ```
   - No reverse dependencies exist (e.g., `inventory` never imports `sale` or `purchase`).

### 3.2. Identified Domain Packages

| Domain Package | Purpose & Scope | Contained Files | Routes Owned |
|---|---|---|---|
| **`internal/auth`** | Authentication, password verification, JWT issuance, profile inspection. | `handler.go`, `service.go`, `routes.go` | `POST /api/v1/auth/login`, `GET /api/v1/auth/me` |
| **`internal/user`** | User lifecycle, registration, role assignment, account deactivation. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST /api/v1/users`, `PATCH /api/v1/users/:id/role` |
| **`internal/category`** | Category catalog management and queries. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST, PUT, DELETE /api/v1/categories` |
| **`internal/product`** | Product catalog, SKU uniqueness, pricing, category filtering. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST, PUT, DELETE /api/v1/products` |
| **`internal/inventory`** | Stock levels, reorder alerts, atomic stock adjustments. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET /api/v1/inventory`, `PATCH /api/v1/inventory/:id/adjust` |
| **`internal/supplier`** | Supplier contact records and status tracking. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST, PUT, DELETE /api/v1/suppliers` |
| **`internal/customer`** | Customer contact records and purchase history. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST, PUT, DELETE /api/v1/customers` |
| **`internal/purchase`** | Purchase orders, line items, supplier order receiving transactions. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST /api/v1/purchases`, `POST /api/v1/purchases/:id/receive` |
| **`internal/sale`** | Sales checkout, price verification, atomic stock decrements, refunds. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET, POST /api/v1/sales`, `POST /api/v1/sales/:id/refund` |
| **`internal/report`** | Aggregated analytics, revenue summaries, sales velocity, low-stock metrics. | `handler.go`, `service.go`, `repository.go`, `routes.go` | `GET /api/v1/reports/summary`, `GET /api/v1/reports/sales` |

### 3.3. Layer Responsibilities within Each Domain
Every domain contains clear internal layering:
- **`handler.go`**: Decodes Gin request payload, invokes Gin binding validator, calls domain service, and outputs response via `internal/utils/response`.
- **`service.go`**: Validates business rules, ensures invariants, coordinates multi-step domain logic, and returns sentinel domain errors.
- **`repository.go`**: Executes GORM queries against PostgreSQL. Contains zero HTTP or authorization logic.
- **`routes.go`**: Registers the domain's endpoints onto the Gin router group with the appropriate `RequireRole` middleware.

---

## 4. Authorization & Access Control

Authentication uses Bearer JWT tokens in the `Authorization` header. Authorization is enforced by `internal/middleware/auth.go` (`JWTAuthMiddleware`) and `internal/middleware/rbac.go` (`RequireRole`).

| Domain / Endpoint | Staff | Manager | Admin |
|---|---|---|---|
| `POST /api/v1/auth/login`, `GET /api/v1/auth/me` | Yes | Yes | Yes |
| `GET /api/v1/products`, `GET /api/v1/categories` | Read | Read | Read/Write |
| `POST /api/v1/products`, `PUT /api/v1/products/:id` | No | Read/Write | Read/Write |
| `GET /api/v1/inventory` | Read | Read | Read |
| `PATCH /api/v1/inventory/:id/adjust` | No | Yes | Yes |
| `GET /api/v1/suppliers`, `POST /api/v1/suppliers` | Read | Read/Write | Read/Write |
| `GET /api/v1/purchases`, `POST /api/v1/purchases` | No | Read/Write | Read/Write |
| `POST /api/v1/purchases/:id/receive` | No | Yes | Yes |
| `GET /api/v1/customers`, `POST /api/v1/customers` | Read/Create | Read/Write | Read/Write |
| `GET /api/v1/sales`, `POST /api/v1/sales` | Read/Create | Read/Create | Read/Create |
| `POST /api/v1/sales/:id/refund` | No | Yes | Yes |
| `GET /api/v1/reports/*` | No | Read | Read |
| `GET /api/v1/users`, `POST /api/v1/users` | No | Read | Read/Write |
| `PATCH /api/v1/users/:id/role` | No | No | Read/Write |

---

## 5. File & Module Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                         # Server entry point, DB init, domain DI, router setup
│   ├── internal/
│   │   ├── config/                             # Environment config loading
│   │   ├── database/
│   │   │   ├── connection.go                   # PostgreSQL connection & pool config
│   │   │   └── migrations/                     # Numbered SQL migration scripts
│   │   ├── middleware/                         # Auth JWT, RBAC RequireRole, Logger, Recovery, CORS
│   │   ├── models/                             # Shared GORM entity structs and domain enums
│   │   │   ├── errors.go                       # Shared domain sentinel errors
│   │   │   └── models.go                       # BaseModel, User, Product, Inventory, Sale, Purchase, etc.
│   │   ├── auth/                               # Auth domain: handler, service, routes
│   │   ├── user/                               # User domain: handler, service, repository, routes
│   │   ├── category/                           # Category domain: handler, service, repository, routes
│   │   ├── product/                            # Product domain: handler, service, repository, routes
│   │   ├── inventory/                          # Inventory domain: handler, service, repository, routes
│   │   ├── supplier/                           # Supplier domain: handler, service, repository, routes
│   │   ├── customer/                           # Customer domain: handler, service, repository, routes
│   │   ├── purchase/                           # Purchase domain: handler, service, repository, routes
│   │   ├── sale/                               # Sale domain: handler, service, repository, routes
│   │   ├── report/                             # Report domain: handler, service, repository, routes
│   │   └── utils/
│   │       ├── auth/                           # Password hashing (bcrypt) and JWT utilities
│   │       ├── logger/                         # zerolog logger and LogError helper
│   │       └── response/                       # Standard response envelope builders
│   ├── docs/                                   # Swaggo generated OpenAPI specs
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/                                # Axios client instance with auth interceptors
│   │   ├── components/                         # UI components (buttons, modals, tables, badges)
│   │   ├── hooks/                              # TanStack Query domain hooks (useProducts, useSales)
│   │   ├── layouts/                            # AppLayout, AuthLayout, Sidebar, Navbar
│   │   ├── routes/                             # TanStack Router route definitions
│   │   ├── stores/                             # Zustand stores (useAuthStore, useUIStore)
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── index.css                           # Tailwind CSS styling
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── docs/
│   └── Roadmap.md                              # Core Build & Technology Learning Track checklist
├── AGENTS.md
├── design.md
├── WORKFLOW.md
├── PRD.md
├── Makefile
└── README.md
```

---

## 6. Concurrency, Atomicity & Idempotency

1. **Atomic Inventory Decrements (Sales)**:
   - When a sale is completed, stock decrement is wrapped in a database transaction with `CHECK (quantity >= 0)`.
   - If available stock is insufficient or concurrent writes deplete stock, the transaction rolls back cleanly with `ErrInsufficientInventory`.
2. **Purchase Receiving Idempotency**:
   - Transitioning an order to `received` checks that current status is `pending`. If already `received`, the operation is rejected, preventing double-counting of inventory.
3. **Server-Side Price Calculation**:
   - The backend recalculates all item subtotals and order totals directly from `products.unit_price` or `purchase_items.unit_cost`, preventing client tampering.

---

## 7. Environments & Observability

- **Environment Configuration**: Loaded from environment variables (`APP_ENV`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`).
- **Structured Logging**: `zerolog` outputs colorized console logs in `development` and JSON logs in `production`.
- **Request Tracing**: `X-Request-ID` is assigned to each request by middleware and included in all log entries.
- **Sanitized Errors**: Errors are logged via `logger.LogError(err, summary, fields)`, redacting secrets and avoiding SQL leaks in production.
- **Testing**: Testify assertions and mocks for unit tests across service and handler layers, with integration tests running against a PostgreSQL test database.
