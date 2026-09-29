# Phase 4: Product Management

> Source: `docs/Roadmap.md` → Phase 4
> Status: Not started
> Depends on: Phase 2 (backend foundation), Phase 3 (auth + RBAC middleware)

**What you'll learn:** how a Gin handler binds and validates a request, how a GORM repository
expresses filters/search with a scoped query, and how a service enforces business rules (SKU
uniqueness, category existence) without touching HTTP or SQL.

Reference: `design.md` §3.2 (domain table), §3.3 (layer responsibilities), §4 (RBAC),
§2.3 (categories/products schema).

## Deliverables

### Category domain (`internal/category`)
- [ ] `repository.go` — `CategoryRepository` with `Create`, `GetByID`, `List`, `Update`, `Delete`
      (soft delete via GORM `Delete`).
- [ ] `service.go` — `CategoryService` wrapping the repository; enforce unique category name
      (translate PostgreSQL unique-violation into a domain sentinel error, e.g. `ErrDuplicate`).
- [ ] `handler.go` — `CategoryHandler` with request structs using Gin `binding` tags; build all
      responses through `internal/utils/response`.
- [ ] `routes.go` — register routes and wire `RequireRole` (writes → `admin`/`manager`).
- [ ] `GET /api/v1/categories` — paginated list.
- [ ] `GET /api/v1/categories/:id`.
- [ ] `POST /api/v1/categories` — validate `name`, `description`.
- [ ] `PUT /api/v1/categories/:id`.
- [ ] `DELETE /api/v1/categories/:id` — reject or guard if products still reference it.

### Product domain (`internal/product`)
- [ ] `repository.go` — `ProductRepository` with CRUD plus a search/filter query.
- [ ] `service.go` — `ProductService`; verify the target category exists before create/update;
      enforce SKU uniqueness; never trust client-supplied price semantics.
- [ ] `handler.go` — `ProductHandler`; represent `unit_price` and `cost_price` as `int64` minor
      units in request/response DTOs.
- [ ] `routes.go` — register routes and wire `RequireRole`.
- [ ] `GET /api/v1/products` — paginated list.
- [ ] `GET /api/v1/products/:id`.
- [ ] `POST /api/v1/products` — `sku`, `name`, `description`, `unit_price`, `cost_price`,
      `category_id`, `is_active`.
- [ ] `PUT /api/v1/products/:id`.
- [ ] `DELETE /api/v1/products/:id` — soft delete.

### Search & filtering
- [ ] Search products by `name` or `sku` using a case-insensitive `ILIKE`/`lower()` query.
- [ ] Filter products by `category_id`.
- [ ] Support pagination (`page`, `per_page`) and return pagination metadata via
      `response.Paginated`.
- [ ] Add appropriate indexes if a new query path needs one (products already index `name`,
      `category_id`, `is_active`).

### Wiring & tests
- [ ] Wire the new repositories/services/handlers into `cmd/server/main.go`.
- [ ] Add table-driven unit tests for services (SKU conflict, missing category, not found).
- [ ] Add handler tests using `httptest` + Gin test mode, asserting response envelopes.
- [ ] Add repository tests against the test database (or mocks consistent with existing patterns).
- [ ] Regenerate Swagger docs (`make swagger`).

## Notes

- Service layer owns business rules; repositories only run queries; handlers only do HTTP.
- A missing category or duplicate SKU is a business error, not a 500 — map it to `400`/`409`.
- Prices are integers in minor units; no `float32`/`float64` anywhere.

## Definition of Done

- [ ] Code compiles cleanly; no cyclic imports.
- [ ] Layer boundaries respected (no GORM in handlers, no Gin context in services).
- [ ] `response.Success` / `response.Paginated` / `response.Error` used for every response.
- [ ] Failures logged via `logger.LogError` with structured fields and no secrets.
- [ ] Role checks applied on write routes.
- [ ] Unit/integration tests pass (`make test`).
- [ ] Swaggo annotations present and `make swagger` re-run.
