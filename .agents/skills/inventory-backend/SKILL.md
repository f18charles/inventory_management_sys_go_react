---
name: inventory-backend
description: >-
  Project-specific conventions for the Inventory Management System's Go backend
  (Gin, GORM, PostgreSQL/Neon, golang-migrate, JWT auth, Testify, zerolog,
  Swagger). Use this skill whenever writing, scaffolding, reviewing, or fixing
  ANY backend code in this repo's backend directory — domain handlers, services,
  repositories, routes, models, migrations, middleware, the router, or tests.
  Also use it when reviewing existing backend code for convention drift (DB
  queries in handlers, business logic in repositories, missing transactions
  around multi-record writes, floating-point money, missing logging, wrong
  package layout), even if the user didn't explicitly ask for a review. Trigger
  on requests like "add a purchases endpoint," "create a service for X," "write
  the repository for Y," "add a migration for Z," or "does this handler look
  right." Encodes the domain-based package layout plus the Handler → Service →
  Repository layering, transaction, error, and logging conventions defined in
  AGENTS.md and design.md.
---

# Inventory Backend Skill

Project-specific backend conventions for the Inventory Management System. Scoped to `backend/` in this repo — not a general Go/Gin skill. Pairs with the `inventory-frontend` skill on the other side of the API boundary.

## Package layout (domain-based, non-negotiable)

Organize by **feature, not by layer**. Each domain is a self-contained package under `internal/<domain>/` containing `handler.go`, `service.go`, `repository.go` (omit repository only when the domain has no persistence), and `routes.go`. This mirrors `design.md` §3.2/§5 and `AGENTS.md` ground rule 1:

```
backend/
├── cmd/server/main.go     # entry point: config, DB, logger, router
└── internal/
    ├── auth/      handler.go, service.go, routes.go
    ├── user/      handler.go, service.go, repository.go, routes.go
    ├── category/  handler.go, service.go, repository.go, routes.go
    ├── product/   handler.go, service.go, repository.go, routes.go
    ├── inventory/ handler.go, service.go, repository.go, routes.go
    ├── supplier/  handler.go, service.go, repository.go, routes.go
    ├── customer/  handler.go, service.go, repository.go, routes.go
    ├── purchase/  handler.go, service.go, repository.go, routes.go
    ├── sale/      handler.go, service.go, repository.go, routes.go
    ├── report/    handler.go, service.go, repository.go, routes.go
    ├── models/    shared GORM entities + sentinel errors (imports no domain)
    ├── middleware/ JWTAuthMiddleware, RequireRole, CORS, logger, recovery
    ├── router/    Gin engine setup + domain DI / route mounting
    ├── config/    env-driven configuration (incl. DATABASE_URL / TEST_DATABASE_URL)
    ├── database/  connection + migrations/*.sql
    └── utils/
        ├── auth/     bcrypt password hashing + JWT issue/parse
        ├── logger/   zerolog wrapper + LogError (dev/prod detail switch)
        └── response/ envelopes + shared RespondError / BindJSON
```

Never create layer packages like `internal/handlers`, `internal/services`, or `internal/repositories` — the layer is a _file_ inside the domain, not a folder.

Cross-domain dependencies form a **DAG** (e.g. `auth -> user`, `sale -> inventory, product, customer`, `purchase -> inventory, product, supplier`); depend on the target domain's exported interface (`user.Repository`, `inventory.Service`) and never import back or create a cycle. `internal/models` is the only shared dependency and imports no domain. See `design.md` §3.1.

## When generating new code

Follow AGENTS.md's feature order for anything beyond a one-line fix:

```
Migration → Repository (+tests) → Service (+tests) → Handler (+tests) → Swagger annotations
```

Decision path for a single new piece of code:

1. **Touches the database schema?** → a new versioned migration pair under `internal/database/migrations/`. Never edit an applied migration; add a new one.
2. **Reads/writes rows?** → a method on the owning domain's repository interface + struct in `internal/<domain>/repository.go`. See `templates/repository.template.go`.
3. **Makes a business decision** (can this happen, is there enough stock, is this allowed)? → a method on the owning domain's service in `internal/<domain>/service.go`, depending on repository *interfaces* (DIP), never a concrete GORM type. See `templates/service.template.go`.
4. **Exposes an HTTP endpoint?** → a Gin handler in `internal/<domain>/handler.go` plus a route registration in `internal/<domain>/routes.go`, and a Swagger annotation block. See `templates/handler.template.go` and `templates/routes.template.go`.
5. **Multiple records must change together** (sale + sale items + inventory; purchase receipt + inventory)? → wrap the repository calls in a `db.Transaction(...)` inside the **service**, never in the handler or repository. See the transaction example in `templates/service.template.go`.

Always check the owning `internal/<domain>/` package and `internal/models/errors.go` for something that already covers the need before adding a new file or a new sentinel error.

## Layering rules (non-negotiable)

```
Gin Handler → Service → Repository → GORM → PostgreSQL
```

- **Handler** (`internal/<domain>/handler.go`): binds/validates the request (Gin binding tags via `response.BindJSON`), extracts auth context, calls exactly one service method, maps the returned error via `response.RespondError`, and logs only handler-specific failures. No `db.` calls. No business logic (no "if quantity > stock" checks here).
- **Service** (`internal/<domain>/service.go`): pure business logic, returns sentinel errors (`models.ErrInsufficientInventory`, not `http.StatusBadRequest`), owns transaction boundaries, logs the meaningful business event and its outcome. Never imports `gin` or references `*gin.Context`.
- **Repository** (`internal/<domain>/repository.go`): pure persistence — CRUD + queries, returns models or wrapped persistence errors, no business decisions, no HTTP, no JWT checks. Each method accepts a `*gorm.DB` (or a transaction handle passed in) so it works inside or outside a service transaction.
- **Routes** (`internal/<domain>/routes.go`): mounts the domain's endpoints under `/api/v1` with the right `JWTAuthMiddleware` / `RequireRole` chain.

If asked to "just quickly" query the database from a handler for convenience, don't — route it through a service/repository call instead and say so.

## Error handling conventions

- Sentinel errors live in `internal/models/errors.go` (shared) or a domain-local `errors.go` for domain-specific ones, and are reused, not redefined per feature:
  ```go
  var (
      ErrNotFound              = errors.New("resource not found")
      ErrUnauthorized          = errors.New("unauthorized")
      ErrForbidden             = errors.New("forbidden")
      ErrInsufficientInventory = errors.New("insufficient inventory")
  )
  ```
- Repositories wrap unexpected persistence errors with context (`fmt.Errorf("get product %s: %w", id, err)`) but return `models.ErrNotFound` (not raw `gorm.ErrRecordNotFound`) and map `gorm.ErrDuplicatedKey` to `models.ErrConflict` (the connection enables `TranslateError: true`), so services never import gorm's error types.
- Handlers translate sentinel errors to HTTP status + envelope through the single shared `response.RespondError` — don't hand-roll status-code `switch` statements in every handler.
- Never let a raw SQL error, GORM error, or stack trace reach the JSON response.

## Transactions

Any service method that changes more than one related record wraps the work in `db.Transaction(func(tx *gorm.DB) error { ... })`, passing `tx` into repository calls. Log the attempt and, critically, the outcome — success with resulting IDs, or failure with which step failed. A transaction that fails without a log line is a bug.

## Money

All money fields are `int64`, smallest currency unit. Never `float64`/`float32` for price, cost, or total fields — flag it immediately if seen. Line totals (`quantity × unit price`) are always computed in the service from trusted server-side data, never accepted as-is from the request body.

## Logging

Per AGENTS.md: `zerolog`, but with a single `APP_ENV`-driven switch (`internal/utils/logger`, see `templates/logger.template.go`) controlling both format and verbosity:

- **Development** (`APP_ENV != production`): colorized console output, `debug` level, and error logs include the raw underlying error text via `logger.LogError`. Use this when actively investigating an issue.
- **Production** (`APP_ENV == production`): plain JSON output, `info` level, and `logger.LogError` logs a safe summary with `detail_suppressed: true` instead of the raw error text. To see the exact underlying issue, restart with `APP_ENV=development` — the switch is the environment variable, not a code change.
- Known-sensitive keys (`password`, `pass_hash`, `token`, `authorization`, `jwt_secret`, `secret`) are stripped from structured fields in **both** modes — the dev/prod switch controls error-detail verbosity, it is never a reason to log credentials.

For significant business events that aren't failures (sale completed, purchase received), log normally with `zerolog`'s regular API — `logger.LogError` is specifically for the failure path. Request-logging middleware handles every request automatically; don't duplicate that in handlers.

## Database & environment

- **Migrations** use `golang-migrate` with numbered SQL files in `internal/database/migrations/` (`000001_*.up.sql` / `.down.sql`). **Never** use GORM `AutoMigrate` for schema definition or changes. Once a migration has successfully applied anywhere, never edit it — add a new versioned migration.
- **Application database** is Neon PostgreSQL via `DATABASE_URL` (TLS; keep `sslmode=require` in the URL). `make migrate-up` and the server target this.
- **Tests** run against the **local** PostgreSQL instance via `TEST_DATABASE_URL` (never Neon). When `TEST_DATABASE_URL` is unset, config falls back to the discrete `DB_*` values.
- `config.Load()` resolves `DSN()` (app) and `TestDSN()` (tests) accordingly; don't bypass it with hard-coded connection strings.

## API & Swagger conventions

- Routes under `/api/v1/...`, registered by the owning domain's `internal/<domain>/routes.go` and mounted by `internal/router` (which performs dependency injection).
- Every handler gets a `swaggo`-style comment block directly above it (method, path, summary, request body, response codes/shapes) — see `templates/handler.template.go`.
- **Swagger generation is not automatic.** Run `make swagger` (`swag init -g cmd/server/main.go -o docs`) before committing handler changes. See `templates/swagger-regen.template.txt`.
- All responses go through `internal/utils/response` — never build the envelope with raw `gin.H{...}` in a handler:
  ```json
  { "data": { "...": "..." } }
  { "error": { "code": "PRODUCT_NOT_FOUND", "message": "Product not found" } }
  { "data": [ ... ], "meta": { "page": 1, "page_size": 20, "total_items": 134, "total_pages": 7 } }
  ```

## Testing conventions

Testify (`assert`/`require`, `suite` where it cuts boilerplate). Tests live under `internal/tests/<domain>/`, with shared testify mocks in `internal/tests/mocks/`:

- **Repository tests** (`internal/tests/<domain>/repository_test.go`) run against the real local test PostgreSQL (`TEST_DATABASE_URL`) — not mocked — since the point is verifying actual SQL/constraints. Wrap each case in a transaction that rolls back, and skip (not fail) when no test database is reachable.
- **Service tests** (`internal/tests/<domain>/service_test.go`) mock the repository *interface* (the concrete payoff of DIP) and assert business outcomes: uniqueness conflicts, inactive accounts, role validation, "can't oversell stock". Methods that open a `db.Transaction` need a real `*gorm.DB`, so cover those with a repository/integration test instead.
- **Handler tests** (`internal/tests/<domain>/handler_test.go`) use `httptest` + Gin's test mode, registering the real domain routes to assert status codes, RBAC, and envelope shape — not internal implementation.
- Table-driven tests for anything with more than two meaningful cases. See `templates/service_test.template.go`.

## Reviewing existing code

When asked to review, or touching a file for an unrelated reason, check for and flag:

- [ ] Code in a layer folder (`internal/handlers|services|repositories`) instead of the owning `internal/<domain>/` package
- [ ] A cross-domain import that reverses or cycles the DAG (e.g. `user` importing `auth`)
- [ ] A GORM/`db.` call inside a handler
- [ ] Business logic (stock checks, permission decisions) inside a repository or handler instead of a service
- [ ] A service returning an `http.Status*` code or importing `gin`
- [ ] A multi-record write (sale+items+inventory, purchase receipt+inventory) not wrapped in `db.Transaction`
- [ ] `float64`/`float32` used for any money field
- [ ] Client-provided totals trusted instead of recomputed server-side
- [ ] A handler building `gin.H{"data": ...}` / `gin.H{"error": ...}` inline instead of calling `internal/utils/response`
- [ ] A raw GORM/SQL error or stack trace reaching the API response
- [ ] A new sentinel error defined redundantly instead of reusing an existing one
- [ ] A transaction or significant business event with no corresponding log line
- [ ] An error logged with `log.Error().Err(err)` directly instead of via `logger.LogError`
- [ ] Sensitive fields (password, token, etc.) passed into a logged fields map
- [ ] A handler or endpoint change with no matching Swagger annotation update, or annotations updated without re-running `make swagger`
- [ ] A repository interface with no corresponding mock used in service tests
- [ ] Tests that should hit `TEST_DATABASE_URL` silently connecting to (or requiring) the application/Neon database

Report findings as a short, concrete list — then apply fixes if asked to fix rather than just review.

## Templates

Read the relevant template before writing the corresponding file — they encode the exact interface shapes, error handling, and logging calls used across this codebase:

- `templates/repository.template.go` — domain repository interface + GORM implementation
- `templates/service.template.go` — domain service with DI, sentinel errors, transaction + logging example
- `templates/handler.template.go` — domain Gin handler using the response package + Swagger annotation
- `templates/routes.template.go` — domain `routes.go` + `internal/router` dependency injection
- `templates/response.template.go` — uniform success/error/pagination helpers plus `RespondError` / `BindJSON`
- `templates/logger.template.go` — dev/production logging switch with sensitive-field redaction
- `templates/swagger-regen.template.txt` — how and when to regenerate Swagger docs
- `templates/service_test.template.go` — Testify service test with a mocked repository
