# Inventory Management System

A full-stack inventory management system for small and medium-sized businesses.

The application manages products, categories, suppliers, customers, stock levels, purchases, sales, users, and inventory activity through a Go backend and React frontend.

The project is being built with a focus on understanding software architecture, database design, backend engineering, testing, and production-oriented development practices — and doubles as a hands-on way to learn tools commonly used in real engineering jobs. See `ROADMAP.md` for the phased build plan and `AGENTS.md` for the engineering rules an AI coding agent (or a human contributor) should follow in this repo.

## Tech Stack

### Backend
- Go
- Gin
- GORM
- PostgreSQL
- golang-migrate
- REST API
- JWT authentication
- bcrypt/Argon2 password hashing
- Testify (testing)
- Swagger/OpenAPI (API documentation)

### Frontend
- React
- Vite
- TanStack Router
- TanStack Query
- Zustand
- axios
- Tailwind CSS

### Infrastructure
- Docker
- Docker Compose
- PostgreSQL
- GitHub Actions
- Environment-based configuration

## Core Features

### 1. Authentication and Authorization

Users can securely access the system according to their assigned role.

Supported roles: **Admin**, **Manager**, **Staff**

Features:
- User registration by authorized users
- Login / Logout
- Password hashing
- JWT-based authentication with token validation
- Role-based authorization
- Activate and deactivate users
- User profile management

Example permissions:

| Action            | Admin | Manager |   Staff |
| ----------------- | ----: | ------: | ------: |
| Manage users      |   Yes |      No |      No |
| Manage products   |   Yes |     Yes | Limited |
| Manage categories |   Yes |     Yes | Limited |
| Manage suppliers  |   Yes |     Yes |      No |
| Manage customers  |   Yes |     Yes |     Yes |
| Create purchases  |   Yes |     Yes |      No |
| Create sales      |   Yes |     Yes |     Yes |
| View reports      |   Yes |     Yes | Limited |
| Manage inventory  |   Yes |     Yes | Limited |

Authorization rules are ultimately enforced by the backend, not only by the frontend.

### 2. User Management

Admins can manage system users. A user contains: ID, first name, last name, username, email, password hash, role, active status, timestamps. The password hash is never returned through the API.

### 3. Category Management

Categories group products (e.g. Electronics, Furniture, Stationery, Clothing, Food). Supports create/read/update/delete, viewing products in a category, and preventing deletion when business rules don't allow it.

### 4. Product Management

A product contains: ID, category, SKU, name, description, selling price, cost price, active status, timestamps. Supports create/read/update, search by SKU/name, filter by category, and viewing product inventory. Products are deactivated rather than permanently deleted when historical sales or purchases reference them.

### 5. Inventory Management

A stock record contains: product, quantity, reorder level, last updated timestamp. Supports viewing current stock, increasing/decreasing stock, low-stock detection, reorder levels, and inventory history, while preventing stock from going negative. Inventory changes are tied to business operations (receiving purchases, completing sales); direct manipulation is restricted to authorized operations.

### 6. Supplier Management

A supplier contains: ID, company name, contact person, phone, email, address, timestamps. Supports create/read/update/delete (or deactivate) and viewing purchase history.

### 7. Customer Management

A customer contains: ID, first name, last name, phone, email, address, timestamps. Supports create/read/update/delete (or deactivate) and viewing purchase history.

### 8. Purchase Management

A purchase contains: ID, supplier, creating user, purchase date, status, total amount, notes, purchase items. Statuses: `pending`, `received`, `cancelled`.

```
Create purchase → Add purchase items → Pending → Receive purchase → Increase inventory → Received
```

Cancelling a purchase that hasn't been received does not modify inventory. Receiving a purchase increases inventory. The purchase and inventory update happen inside a database transaction.

### 9. Sales Management

A sale contains: ID, customer, creating user, sale date, status, total amount, payment method, sale items. Statuses: `pending`, `completed`, `cancelled`, `refunded`. Payment methods: `cash`, `card`, `mobile_money`, `bank_transfer`.

```
Create sale → Add sale items → Check inventory → Confirm sale → Decrease inventory → Completed
```

A sale cannot complete without sufficient stock. Creating a completed sale and reducing inventory happen inside one transaction.

### 10. Purchase and Sale Items

Purchase items: purchase ID, product ID, quantity, unit cost, subtotal. Sale items: sale ID, product ID, quantity, unit price, subtotal. Subtotal is always `quantity × unit price`, calculated on the backend — client-provided totals are never trusted.

### 11. Dashboard

Quick view of the business: total/active products, total categories/customers/suppliers, current inventory value, low-stock products, today's and monthly sales/purchases, sales revenue, estimated profit.

### 12. Reports

Sales report, purchase report, inventory report, low-stock report, product performance, customer/supplier purchase history, profit report, sales by payment method, sales by date range — filterable by date, product, category, customer, supplier, and user.

## Database

Core tables: `users`, `suppliers`, `customers`, `categories`, `products`, `inventory`, `purchases`, `purchase_items`, `sales`, `sale_items`.

```
Category → Products → Inventory
Supplier → Purchases → Purchase Items → Products
Customer → Sales → Sale Items → Products
User → Purchases, Sales
```

Integrity is enforced through primary keys, foreign keys, unique constraints, not-null constraints, check constraints, indexes, and transactions — the database is the final authority over data integrity, not just the service layer.

Migrations are versioned SQL files managed by `golang-migrate` (never GORM `AutoMigrate`). See `AGENTS.md` for full migration conventions.

## HTTP API

REST API under `/api/v1/...`, built with Gin. Example endpoints:

```
/api/v1/auth/login
/api/v1/users, /api/v1/users/{id}
/api/v1/categories, /api/v1/categories/{id}
/api/v1/products, /api/v1/products/{id}
/api/v1/inventory, /api/v1/inventory/{productId}
/api/v1/suppliers, /api/v1/suppliers/{id}
/api/v1/customers, /api/v1/customers/{id}
/api/v1/purchases, /api/v1/purchases/{id}
/api/v1/sales, /api/v1/sales/{id}
```

Consistent response envelopes:

```json
{ "data": { "id": "..." } }
```
```json
{ "error": { "code": "PRODUCT_NOT_FOUND", "message": "Product not found" } }
```

Interactive API documentation is generated via Swagger/OpenAPI and served at `/swagger/index.html` in development.

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   ├── handlers/
│   ├── middleware/
│   ├── validators/
│   └── utils/
├── migrations/
├── tests/
├── .env
├── .gitignore
├── go.mod
└── go.sum

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   ├── layouts/
│   ├── hooks/
│   ├── stores/        # Zustand stores
│   ├── routes/         # TanStack Router route tree
│   ├── api/
│   ├── utils/
│   └── types/
├── public/
├── package.json
└── vite.config.js
```

## Development Environment

PostgreSQL (and, later, Redis/Kafka/BI tooling — see `ROADMAP.md`) run through Docker Compose during development. The React app can run directly on the host during early development; the full stack may eventually run entirely through Docker Compose.

### Backend Setup
```bash
cd backend
go mod download
go run cmd/api/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## Environment Variables

```
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=inventory
DB_USER=postgres
DB_PASSWORD=

JWT_SECRET=
JWT_EXPIRATION=
```

Production secrets must never be committed to Git. `.env` is git-ignored; a `.env.example` is provided.

## External Services & Manual Setup

As part of the technology learning track (`ROADMAP.md`), a few tools require an account, subscription, or cloud console setup that has to be done by a person — an AI coding agent can wire up the integration code, but it can't create accounts, accept terms of service, or hold billing/credentials on your behalf. Do these manually when you reach the relevant roadmap step, then hand the resulting keys/config to the app via environment variables (never commit them):

- **Sentry** — create a Sentry account/project (backend and frontend), generate the DSNs, and add them as env vars (`SENTRY_DSN_BACKEND`, `SENTRY_DSN_FRONTEND` or similar). Once set, the agent can implement the SDK integration and error-reporting middleware.
- **PostHog** — create a PostHog account/project (or self-host), generate an API key, and add it as an env var. Once set, the agent can implement the frontend tracking calls and event schema.
- **AWS** — create/access an AWS account, set up billing, and provision or authorize the chosen services (RDS, S3, ECS/EC2, IAM roles, etc.) through the console or your own IaC credentials. Once access/credentials exist, the agent can write the infrastructure config and deployment code that targets them.

Redis, Kafka, and BI tooling (e.g. Metabase) are **not** in this list — they're run locally via Docker Compose and can be fully set up in code, so they stay on the normal roadmap as agent-buildable work.

## Security

- Never store plaintext passwords or return password hashes in API responses.
- Validate authentication tokens; enforce authorization on the backend.
- Validate all user input; use parameterized queries through GORM.
- Store secrets in environment variables; never commit `.env`.
- Avoid exposing database errors to clients; apply reasonable request limits; use HTTPS in production.

## Testing

Testing exists at multiple levels — model, repository, service, handler, and integration — using Testify. See `AGENTS.md` for what each level should cover. The goal is confidence in important behavior, not maximum test count.

## Learning Goals

This project is also a practical software engineering learning project. Core concepts being practiced: abstraction, composition, separation of concerns (SOLID), domain modeling, database normalization, referential integrity, SQL, database migrations, ORM design, the repository/service pattern, REST API design, authentication/authorization, transactions, indexing, validation, testing, dependency management, configuration management, Docker, and CI/CD — plus, per the technology learning track in `ROADMAP.md`, hands-on exposure to caching (Redis), event streaming (Kafka), observability (Sentry), product analytics (PostHog), BI tooling, and cloud infrastructure (AWS).

The goal is to understand why each layer and tool exists, not simply to make the application work.

## License

This project is currently intended as a personal software engineering and portfolio project. License information can be added when the project is ready for public distribution.