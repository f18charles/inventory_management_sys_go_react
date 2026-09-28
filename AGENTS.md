# Instructions for AI Agents Building This Repository

Read this document together with `PRD.md` (what to build and product scope), `design.md` (architecture, schema, and domain-based package layout), and `WORKFLOW.md` (git flow, milestones, and issue tracking). Do not re-derive architecture that is already decided in `design.md`. Follow it strictly, and if a decision looks incorrect, flag it explicitly rather than silently diverging.

Before making any change, and again before committing or pushing, read `WORKFLOW.md`. It defines branch naming, commit message formatting, and PR generation. Every commit references the GitHub issue it belongs to.

---

## 1. Ground Rules

1. **Follow the Domain Layered Architecture**: Each domain under `internal/<domain>/` (`auth`, `user`, `category`, `product`, `inventory`, `supplier`, `customer`, `purchase`, `sale`, `report`) strictly follows:
   ```
   Handler (handler.go) → Service (service.go) → Repository (repository.go) → GORM → PostgreSQL
   ```
   Never bypass layers. Handlers do not execute SQL or import GORM. Repositories do not perform business decisions or inspect JWT tokens. Services do not touch Gin contexts or build HTTP status codes (see `design.md` §3).

2. **Prevent Cyclic Imports via Shared Models & DAG**: All GORM structs live in `internal/models`. Domain packages import `internal/models`. Cross-domain dependencies follow a strict directed acyclic graph (e.g. `sale` depends on `inventory`, but `inventory` never imports `sale`). Do not create circular package dependencies (see `design.md` §3.1).

3. **Always Execute Multi-Record Writes in Transactions**: Any operation modifying inventory or multiple tables (`Sale` + `SaleItems` + `Inventory`, or `Purchase` receiving + `Inventory`) must run inside an explicit database transaction with rollback on failure (see `design.md` §3.2 and §6).

4. **Store Currency in Minor Units as Integers**: Never use floating-point types (`float32`, `float64`, `numeric` with decimals) for currency. Store all monetary values as `int64` (cents). Always recompute prices, line item subtotals, and order totals server-side (see `design.md` §2 and §6).

5. **Zero Negative Stock**: Stock counts must never become negative. Verify available inventory before decrementing, and rely on the database check constraint `CHECK (quantity >= 0)` as the final safeguard (see `design.md` §2.3 and §6).

6. **Sanitize Log Output via Shared Helper**: Always log service and handler failures through `logger.LogError(err, safeSummary, fields)` in `internal/utils/logger`. Do not make raw zerolog error calls that leak raw SQL statements or stack traces in production (see `design.md` §7).

7. **Redact Sensitive Credentials**: Never log passwords, password hashes (`pass_hash`), JWT secrets, or authorization headers. This is enforced globally by the redaction filter in `internal/utils/logger`.

8. **Standard API Response Envelopes**: Handlers must build responses exclusively using `internal/utils/response` helpers (`response.Success`, `response.Paginated`, `response.Error`). Never construct ad-hoc `gin.H{...}` response maps (see `design.md` §3.3).

9. **Migrations via golang-migrate Only**: Never use GORM `AutoMigrate` for schema definition or schema changes. All database updates use numbered SQL files in `backend/internal/database/migrations/` (`000001_*.up.sql` and `000001_*.down.sql`). Once applied, never edit an existing migration (see `design.md` §1 and §2.3).

10. **Backend Is the Authoritative Auth Boundary**: Never rely on frontend UI hiding or disabling for access control. Enforce role checks (`admin`, `manager`, `staff`) on API endpoints using `RequireRole` middleware in `internal/middleware` (see `design.md` §4).

11. **Respect Scope Limits from PRD**: Do not build multi-tenant organization switches, external automated payment webhooks, or email delivery queues in v1 (see `PRD.md` §6).

---

## 2. Conventions

- **Language & Go Style**: Idiomatic Go 1.22+. Explicit error handling with `if err != nil { return ... }`. Small, focused functions. Dependency injection via struct constructors (`NewRepository`, `NewService`, `NewHandler`).
- **Domain Organization (Package by Feature)**: Group related capabilities into domain packages under `internal/<domain>` containing `handler.go`, `service.go`, `repository.go`, and `routes.go`.
- **Validation**: Request structural validation lives in Gin request structs using `binding:"required,min=1"` tags in `handler.go`. Business validation (stock availability, role permissions, state transitions) lives in `service.go`. Database legality lives in PostgreSQL constraints.
- **Naming**:
  - Database tables: snake_case plural (`users`, `purchase_items`, `sale_items`, `inventories`).
  - Database columns: snake_case (`pass_hash`, `company_name`, `unit_price`, `reorder_level`, `subtotal`, `sub_total`).
  - Go types: PascalCase (`User`, `Category`, `Product`, `Inventory`, `Purchase`, `Sale`).
  - React components: PascalCase (`ProductTable.jsx`, `SaleModal.jsx`).
  - React hooks and stores: camelCase (`useProducts.js`, `useAuthStore.js`).
- **File Placement**: Follow the domain tree in `design.md` §5 strictly.
- **Error Handling**: Define sentinel errors in `internal/models/errors.go` or domain packages (e.g. `ErrNotFound`, `ErrInsufficientInventory`, `ErrInvalidStatusTransition`). Handlers translate domain errors to proper HTTP status codes (`400`, `401`, `403`, `404`, `409`, `422`, `500`).
- **API Documentation**: Annotate all Gin handler functions with Swaggo comments. Run `swag init -g cmd/server/main.go -o docs` whenever endpoint annotations or request/response structs change.
- **Frontend Architecture**: React with plain JavaScript/JSX. TanStack Router for route tree in `src/routes`. TanStack Query for server data fetching in `src/hooks`. Zustand in `src/stores` for authentication session and local UI state. Axios in `src/api/client.js` for API requests.

---

## 3. Definition of Done for a Feature

Before considering any backend or frontend feature complete, verify every item on this checklist:

- [ ] Code compiles cleanly without errors or warnings.
- [ ] No cyclic import errors across domain packages.
- [ ] Database migrations have matching `.up.sql` and `.down.sql` scripts in `backend/internal/database/migrations`.
- [ ] Layering boundaries are respected (no database queries in handlers; no Gin context in services; no business logic in repositories).
- [ ] Multi-record mutations execute within an atomic database transaction.
- [ ] Monetary amounts are integers in minor units; line item totals and sale sums are recalculated on the backend.
- [ ] Unit and integration tests pass using Testify (`assert`/`require`).
- [ ] Handler responses use standard `internal/utils/response` envelopes (`data`, `error`, `meta`).
- [ ] Failures are logged with `logger.LogError` with appropriate structured fields.
- [ ] Swaggo annotations are present on all modified or new endpoints, and `swag init` was re-run.
- [ ] Role authorization checks are applied to protected routes.

---

## 4. What NOT to Do

- Do not create circular package dependencies between domain folders (e.g. `inventory` importing `sale` while `sale` imports `inventory`).
- Do not write raw SQL queries or invoke GORM methods inside Gin handlers.
- Do not pass `*gin.Context` into service or repository methods.
- Do not store or calculate monetary amounts with floating-point types (`float32`, `float64`).
- Do not trust client-supplied item prices or cart totals; always recalculate them from database product records.
- Do not delete or edit already-applied migrations; add a new versioned migration instead.
- Do not expose database errors, connection strings, or internal stack traces in API response payloads.
- Do not use `fmt.Println` or raw un-redacted `log.Print` in production code; use structured `zerolog` via `internal/utils/logger`.
- Do not build deferred v2 features (multi-tenancy, background email workers, external payment gateways) during v1 implementation.