# Product Requirements Document (PRD) — Inventory Management System

Read this document together with `design.md` (technical architecture and database schema), `AGENTS.md` (development ground rules), and `WORKFLOW.md` (milestones, issues, and git conventions).

---

## 1. Executive Summary & Vision

Small and medium-sized businesses (SMBs) struggle with fragmented stock management, manual record-keeping, overselling, and inventory reconciliation errors. This project delivers a high-reliability, full-stack Inventory Management System designed for precision, speed, and strict data consistency.

The system handles product catalogs, category hierarchies, supplier purchases, customer sales, real-time stock levels, role-based access control, and comprehensive audit history. It serves as both a commercial-grade business management tool and a production-grade engineering reference architecture featuring Go, Gin, GORM, PostgreSQL, React, TanStack Router, TanStack Query, and Zustand.

Beyond core inventory workflows, the platform provides a progressive learning and scaling pathway covering caching (Redis), distributed event streaming (Kafka), application monitoring (Sentry), product analytics (PostHog), business intelligence (Metabase), and cloud deployment (AWS).

---

## 2. Target Users & Personas

| Role | Responsibilities | Key Needs in the System |
|---|---|---|
| **Admin** | System administration, user onboarding, business configuration, full financial and stock oversight. | Role management, audit logs, unrestricted access to sales, purchases, catalog, and reporting. |
| **Manager** | Day-to-day operations, purchase order approvals, inventory adjustments, supplier and customer relationship management. | Reorder alerts, supplier tracking, purchase receipt verification, sales monitoring, category administration. |
| **Staff / Cashier** | Point-of-sale operations, order lookup, stock inquiries, customer order processing. | Fast product search, reliable stock validation during checkout, simple receipt and order entry. |

---

## 3. Core Problems & Value Proposition

1. **Elimination of Negative Stock & Overselling**: Inventory transactions execute inside strict database transactions with validation checks. A sale cannot complete if requested quantities exceed available stock.
2. **Financial Precision**: All monetary values are represented as integers in the smallest currency unit (cents or minor units). Floating-point calculation errors are completely prevented.
3. **Data Integrity & Traceability**: Stock cannot change arbitrarily. Every inventory adjustment is tied directly to a purchase receipt, completed sale, cancellation, or authorized manual adjustment record.
4. **Separation of Concerns**: Business logic is strictly isolated in the service layer, keeping handlers focused on HTTP transport and repositories focused on persistence.

---

## 4. Feature Requirements (v1 Scope)

### 4.1. Authentication & Role-Based Access Control (RBAC)
- **User Authentication**: Secure login via username or email with passwords hashed using `bcrypt`.
- **JWT Session Tokens**: Stateless authentication using signed JSON Web Tokens (`jwt/v5`) containing user ID, role, and expiration timestamp.
- **Role Enforcement**: Protected API endpoints enforce role checks (`Admin`, `Manager`, `Staff`).
- **User Management**: Admin users can create users, deactivate accounts, and update assigned roles.

### 4.2. Category & Product Catalog
- **Category Hierarchy**: Organize products into structured categories.
- **Product Management**: SKU (unique identifier), product name, description, unit price (in minor units), cost price, status (active/inactive), and category assignment.
- **Search & Filtering**: Search products by SKU or name; filter products by category or stock availability.

### 4.3. Inventory & Stock Control
- **Per-Product Inventory**: Current stock count, low-stock threshold, and reorder point.
- **Low Stock Alerts**: Real-time identification of items at or below their defined threshold.
- **Atomic Stock Adjustments**: Stock increases on purchase receipt; stock decreases on sale completion. All stock changes run inside atomic database transactions.
- **Negative Stock Prohibition**: Database and service layer constraints enforce `stock >= 0` unconditionally.

### 4.4. Supplier & Purchase Management
- **Supplier Records**: Supplier name, contact email, phone, address, and status.
- **Purchase Orders**: Order workflow moving through defined states: `pending` -> `received` or `cancelled`.
- **Line Items**: Multi-item purchase orders containing quantity and unit cost per product.
- **Receiving Workflow**: Marking a purchase order as `received` increases stock for all line items atomically inside a database transaction. Cancelling an unreceived purchase order leaves stock untouched.

### 4.5. Customer & Sales Management
- **Customer Records**: Customer name, email, phone, address, and status.
- **Sales Transactions**: Sales workflow moving through states: `pending` -> `completed`, `cancelled`, or `refunded`.
- **Payment Methods**: Support for `cash`, `card`, `mobile_money`, and `bank_transfer`.
- **Checkout Validation**: Line item prices are recalculated server-side (`quantity * unit_price`). Available stock is verified before finalizing the transaction.
- **Atomic Completion**: Decrements inventory and creates the sale record within an atomic transaction.

### 4.6. Reporting & Analytics
- **Summary Metrics**: Total products, low-stock count, total revenue, total purchases, and net profit.
- **Sales Velocity**: Top-selling products and category breakdowns over specified date ranges.

### 4.7. Frontend User Interface
- **Modern Dashboard**: High-level overview of revenue, transactions, inventory alerts, and quick actions.
- **Catalog & Inventory Views**: Paginated tables with search, category filtering, and modal forms for adding/editing items.
- **POS / Sales Interface**: Streamlined order creation workflow with instant stock checking and payment method selection.
- **Purchase Order Workspace**: Order builder and one-click receiving actions.

---

## 5. Technology Learning Track & Future Architecture (v2+)

The platform is designed to scale beyond the initial core deployment. Once the core build is stable, the following technologies will be integrated sequentially:

1. **Redis Caching & Rate Limiting**:
   - Cache frequent, read-heavy queries (product catalog, categories).
   - Redis-backed rate limiting middleware to prevent API abuse.
   - Token blocklist for immediate session revocation upon logout.
2. **Kafka Event Streaming**:
   - Publish asynchronous domain events (`SaleCompleted`, `PurchaseReceived`, `StockAdjusted`).
   - Consumer microservice for building immutable audit logs and historical analytics.
3. **Sentry Observability**:
   - Centralized backend and frontend exception tracking with error context and breadcrumbs (with sensitive PII redacted).
4. **PostHog Product Analytics**:
   - User behavior tracking, conversion funnel measurement, and feature adoption analysis.
5. **Business Intelligence (Metabase)**:
   - Dedicated read-only replica connection for Metabase BI dashboards, custom SQL queries, and automated scheduled reports.
6. **AWS Cloud Deployment**:
   - Production deployment utilizing AWS RDS (PostgreSQL), ECS/Fargate container hosting, and AWS Secrets Manager.

---

## 6. Non-Goals (Explicitly Out of Scope for v1)

- **Multi-Tenant SaaS / Organization Switching**: Single-tenant deployment per instance in v1. No multi-tenant data partitioning.
- **Custom Hardware / Barcode Scanner SDKs**: Standard keyboard input emulation for barcode scanners is supported; proprietary USB/Bluetooth driver integrations are excluded.
- **Automated External Payment Gateway Webhooks**: Payment statuses are recorded manually by staff upon payment confirmation at point-of-sale in v1.
- **In-App Email Delivery Engine**: Purchase orders and receipts are rendered on screen and printable; automated transactional email queues are deferred to v2.

---

## 7. Success Criteria & Quality Standards

- **Zero Inconsistent Stock States**: Zero instances of negative stock or orphaned purchase/sale items under concurrent load.
- **Sub-100ms API Latency**: P95 latency under 100ms for read endpoints; P95 under 200ms for multi-item transactional writes.
- **Test Coverage**: Minimum 80% unit and integration test coverage across service, repository, and middleware layers using Testify.
- **API Documentation**: 100% Swagger/OpenAPI documentation coverage for all `/api/v1` routes, verified against committed docs in CI.
