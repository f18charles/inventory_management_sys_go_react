# WORKFLOW.md: Git Flow, Milestones, and Issues

Read this document before making any code change, and again before committing or pushing. It defines how work is branched, committed, and closed out in this repository. It complements `PRD.md` (what to build), `design.md` (architecture and schema decisions), and `AGENTS.md` (development ground rules).

Milestones and issues below define the planned delivery order. Every branch, commit, and pull request references the issue number it belongs to.

---

## 1. Branching Strategy

- `main` is always stable and deployable. No direct commits to `main`.
- Create one branch per issue, branched from `main`:
  ```
  <type>/<issue-number>-<short-slug>
  ```
  Example: `feat/3-jwt-authentication`, `fix/7-prevent-negative-stock`, `test/12-sale-service-tests`.
- `<type>` matches the commit type table in §2.
- Keep each branch strictly scoped to its assigned issue. If extra scope is discovered, open a new issue and branch separately.

---

## 2. Commit Message Format

Follow Conventional Commits with an explicit issue reference:

```
<type>(<scope>): <short summary> (#<issue-number>)
```

| Type | Purpose |
|---|---|
| `feat` | New feature or endpoint |
| `fix` | Bug fix or calculation error correction |
| `chore` | Dependency updates, tooling, or configuration changes |
| `docs` | Documentation or Swagger annotation updates |
| `refactor` | Code restructuring without behavior or API changes |
| `test` | Adding, updating, or fixing tests |
| `style` | Formatting or lint cleanups |

### Commit Rules
- One logical change per commit.
- Use the imperative mood in the summary ("add user repository", not "added user repository").
- Include the issue number in every commit: `feat(auth): add password hashing utility (#3)`.
- Use the commit body to explain the architectural rationale for non-obvious choices.

---

## 3. Pull Request Format & Closing Keywords

Before pushing a finished branch, organize its commits chronologically and prepare the PR description block:

```markdown
## Issue #<N>: <Issue Title>

- <type>(<scope>): <summary> (#<N>)
- <type>(<scope>): <summary> (#<N>)

---
Closes #<N>

<1-2 sentences summarizing the change and confirming passing tests.>
```

Use `Closes #<N>` for feature PRs and `Fixes #<N>` for bug fix PRs.

---

## 4. Pre-Commit Checklist

- [ ] Branch name matches `<type>/<issue-number>-<short-slug>`.
- [ ] Code compiles without errors.
- [ ] Tests covering modified logic pass.
- [ ] Every commit message follows the format in §2.
- [ ] The `AGENTS.md` Definition of Done items pass for all touched files.
- [ ] The PR description block is ready with the closing keyword.

---

## 5. Project Milestones & Issue Breakdown

### Milestone 1: Database Foundation & Schema Setup (Completed)
- **Issue #1**: Set up PostgreSQL connection, GORM configuration, and golang-migrate pipeline.
- **Issue #2**: Create migrations for users, categories, suppliers, customers, products, inventory, purchases, purchase_items, sales, and sale_items with foreign keys and check constraints.

### Milestone 2: Backend Architecture & Core Middleware (Completed)
- **Issue #3**: Initialize Gin router and `/api/v1` route groups with config loading.
- **Issue #4**: Implement standard response envelope helpers (`internal/utils/response`) and error handling middleware.
- **Issue #5**: Configure structured logging (`zerolog`), request ID tracking, and panic recovery middleware.

### Milestone 3: Authentication & User Management (Current / Phase 3)
- **Issue #6**: Implement user repository with CRUD operations and unit tests.
- **Issue #7**: Build password hashing (`bcrypt`) and JWT token utilities (`jwt/v5`).
- **Issue #8**: Implement `AuthService` (login, credentials verification) and `UserService` (user creation, role updates).
- **Issue #9**: Build `JWTAuthMiddleware` and `RequireRole` authorization middleware.
- **Issue #10**: Implement `POST /api/v1/auth/login`, `GET /api/v1/auth/me`, and `/api/v1/users` endpoints with handler tests.

### Milestone 4: Product & Category Management
- **Issue #11**: Implement Category repository, service, and handlers (`GET`, `POST`, `PUT`, `DELETE /api/v1/categories`).
- **Issue #12**: Implement Product repository, service, and handlers (`GET`, `POST`, `PUT`, `DELETE /api/v1/products`) with SKU uniqueness checks.
- **Issue #13**: Add product search by name/SKU and filtering by category.

### Milestone 5: Inventory Control & Low Stock Tracking
- **Issue #14**: Implement Inventory repository and stock adjustment service.
- **Issue #15**: Build low-stock detection query and endpoint (`GET /api/v1/inventory/low-stock`).
- **Issue #16**: Implement manual stock adjustment endpoint (`PATCH /api/v1/inventory/:id/adjust`) with role checks and audit logging.

### Milestone 6: Supplier & Purchase Order Management
- **Issue #17**: Implement Supplier repository, service, and CRUD endpoints.
- **Issue #18**: Implement Purchase order creation service with line items and cost totals.
- **Issue #19**: Build purchase receiving transaction workflow (incrementing inventory atomically and updating status to `received`).

### Milestone 7: Customer & Sales Order Management
- **Issue #20**: Implement Customer repository, service, and CRUD endpoints.
- **Issue #21**: Implement Sale creation service with server-side price recalculation and stock availability validation.
- **Issue #22**: Build sale completion transaction workflow (decrementing inventory atomically and persisting sale records).
- **Issue #23**: Implement sale refund and cancellation workflow.

### Milestone 8: API Documentation & Developer Experience
- **Issue #24**: Add Swaggo annotations across all `/api/v1` handlers and generate OpenAPI specifications in `backend/docs`.
- **Issue #25**: Wire Swagger UI endpoint (`/swagger/index.html`) in development mode and add Makefile target `make swagger`.

### Milestone 9: Frontend Architecture & User Interfaces
- **Issue #26**: Set up Vite React application with Tailwind CSS, TanStack Router, and TanStack Query.
- **Issue #27**: Implement authentication flow, login screen, and Zustand session store (`useAuthStore`).
- **Issue #28**: Build main application layout with navigation sidebar, header, and role-based route guards.
- **Issue #29**: Build Product and Category catalog pages with search, filtering, and modal forms.
- **Issue #30**: Build Inventory monitoring table with low-stock badges and quick adjustment dialog.
- **Issue #31**: Build Point of Sale (POS) checkout interface with instant stock validation.
- **Issue #32**: Build Purchase order creation and receiving interface.
- **Issue #33**: Build Dashboard overview with KPI summary cards and sales velocity charts.

### Milestone 10: Production Readiness & Quality Assurance
- **Issue #34**: Write end-to-end integration test suite against PostgreSQL test database.
- **Issue #35**: Set up GitHub Actions CI pipeline running Go tests, linting, Swagger validation, and frontend build.
- **Issue #36**: Create multi-stage production Dockerfile and Docker Compose deployment configuration.

### Milestone 11: Technology Learning Track
- **Issue #37**: Step A - Integrate Redis caching for product catalog and rate-limiting middleware.
- **Issue #38**: Step B - Implement Kafka event producer for sales/inventory domain events and build audit consumer.
- **Issue #39**: Step C - Integrate Sentry error tracking with PII sanitization.
- **Issue #40**: Step D - Integrate PostHog product analytics in frontend.
- **Issue #41**: Step E - Stand up Metabase BI container connected to PostgreSQL read replica.
- **Issue #42**: Step F - Author AWS deployment manifests (RDS, ECS, Secrets Manager).
