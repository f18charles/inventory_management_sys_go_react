---
name: inventory-backend
description: >-
  Project-specific conventions for the Inventory Management System's Go backend
  (Gin, GORM, PostgreSQL, golang-migrate, JWT auth, Testify, zerolog, Swagger).
  Use this skill whenever writing, scaffolding, reviewing, or fixing ANY backend
  code in this repo's backend directory — handlers, services, repositories,
  models, migrations, middleware, or tests. Also use it when reviewing existing
  backend code for convention drift (DB queries in handlers, business logic in
  repositories, missing transactions around multi-record writes, floating-point
  money, missing logging), even if the user didn't explicitly ask for a review.
  Trigger on requests like "add a purchases endpoint," "create a service for X,"
  "write the repository for Y," "add a migration for Z," or "does this handler
  look right." Encodes the Handler → Service → Repository layering, transaction,
  error, and logging conventions defined in AGENTS.md.
---

# Inventory Backend Skill

Project-specific backend conventions for the Inventory Management System. Scoped to `backend/` in this repo — not a general Go/Gin skill. Pairs with the `inventory-frontend` skill on the other side of the API boundary.

## When generating new code

Follow AGENTS.md's Feature Implementation Order for anything beyond a one-line fix:

```
Migration → Repository (+tests) → Service (+tests) → Handler (+tests) → Swagger annotations
```

Decision path for a single new piece of code:

1. **Touches the database schema?** → a new versioned migration pair under `migrations/`. Never edit an applied migration; add a new one.
2. **Reads/writes rows?** → a method on the relevant repository interface + struct under `internal/repositories/`. See `templates/repository.template.go`.
3. **Makes a business decision** (can this happen, is there enough stock, is this allowed)? → a method on the relevant service under `internal/services/`, depending on repository *interfaces* (DIP), never a concrete GORM type. See `templates/service.template.go`.
4. **Exposes an HTTP endpoint?** → a Gin handler under `internal/handlers/`, plus a route registration and a Swagger annotation block. See `templates/handler.template.go`.
5. **Multiple records must change together** (sale + sale items + inventory; purchase receipt + inventory)? → wrap the repository calls in a `db.Transaction(...)` inside the **service**, never in the handler or repository. See the transaction example in `templates/service.template.go`.

Always check `internal/repositories/`, `internal/services/`, and `internal/models/errors.go` (or equivalent) for something that already covers the need before adding a new file or a new sentinel error.

## Layering rules (non-negotiable)

```
Gin Handler → Service → Repository → GORM → PostgreSQL
```

- **Handler**: binds/validates the request (Gin binding tags), extracts auth context, calls exactly one service method, maps the returned error to an HTTP status + envelope, and logs only handler-specific failures. No `db.` calls. No business logic (no "if quantity > stock" checks here).
- **Service**: pure business logic, returns sentinel errors (`ErrInsufficientInventory`, not `http.StatusBadRequest`), owns transaction boundaries, logs the meaningful business event and its outcome. Never imports `gin` or references `*gin.Context`.
- **Repository**: pure persistence — CRUD + queries, returns models or wrapped persistence errors, no business decisions, no HTTP, no JWT checks. Depends on `*gorm.DB` (or a transaction handle passed in), not on services.

If asked to "just quickly" query the database from a handler for convenience, don't — route it through a service/repository call instead and say so.

## Error handling conventions

- Sentinel errors live in one place (e.g. `internal/models/errors.go` or `internal/services/errors.go`) and are reused, not redefined per feature:
  ```go
  var (
      ErrNotFound              = errors.New("resource not found")
      ErrUnauthorized          = errors.New("unauthorized")
      ErrForbidden             = errors.New("forbidden")
      ErrInsufficientInventory = errors.New("insufficient inventory")
  )
  ```
- Repositories wrap unexpected persistence errors with context (`fmt.Errorf("get product %s: %w", id, err)`) but return `ErrNotFound` (not a raw GORM `ErrRecordNotFound`) so services never import `gorm`'s error types.
- Handlers translate sentinel errors to HTTP status + the standard envelope via a single shared mapping function — don't hand-roll status-code `switch` statements in every handler. See `templates/handler.template.go`.
- Never let a raw SQL error, GORM error, or stack trace reach the JSON response.

## Transactions

Any service method that changes more than one related record wraps the work in `db.Transaction(func(tx *gorm.DB) error { ... })`, passing `tx` into repository calls (repositories should accept a `*gorm.DB` so they work identically inside or outside a transaction). Log the attempt and, critically, the outcome — success with resulting IDs, or failure with which step failed — per AGENTS.md's Logging section. A transaction that fails without a log line is a bug.

## Money

All money fields are `int64`, smallest currency unit. Never `float64`/`float32` for price, cost, or total fields — flag it immediately if seen. Line totals (`quantity × unit price`) are always computed in the service from trusted server-side data, never accepted as-is from the request body.

## Logging

Per AGENTS.md: `zerolog`, but with a single `APP_ENV`-driven switch (`internal/utils/logger`, see `templates/logger.template.go`) controlling both format and verbosity:

- **Development** (`APP_ENV != production`): colorized console output, `debug` level, and — this is the important part — error logs include the raw underlying error text via `logger.LogError`. Use this when actively investigating an issue.
- **Production** (`APP_ENV == production`): plain JSON output, `info` level, and `logger.LogError` logs a safe summary message with `detail_suppressed: true` instead of the raw error text. To see the exact underlying issue, restart with `APP_ENV=development` — the switch is the environment variable, not a code change.
- Known-sensitive keys (`password`, `pass_hash`, `token`, `authorization`, `jwt_secret`, `secret`) are stripped from structured fields in **both** modes — the dev/prod switch controls error-detail verbosity, it is never a reason to log credentials.

For significant business events that aren't failures (sale completed, purchase received), log normally with `zerolog`'s regular API — `logger.LogError` is specifically for the failure path. Request-logging middleware still handles every request automatically; don't duplicate that in individual handlers. See `templates/service.template.go` for both patterns in one workflow.

## API & Swagger conventions

- Routes under `/api/v1/...`, grouped by resource in `internal/handlers/routes.go` (or equivalent router setup).
- Every handler gets a `swaggo`-style comment block directly above it (method, path, summary, request body, response codes/shapes) — see `templates/handler.template.go`.
- **Swagger generation is not automatic.** `swag init` only regenerates `docs/` when you explicitly run it — editing an annotation comment does nothing on its own. Run it via a Makefile target (`make swagger`) before committing handler changes, and add a CI check that fails the build if generated docs are stale. See `templates/swagger-regen.template.txt`.
- All responses go through `internal/utils/response` (see `templates/response.template.go`) — never build the envelope with raw `gin.H{...}` in a handler:
  ```json
  { "data": { "...": "..." } }
  { "error": { "code": "PRODUCT_NOT_FOUND", "message": "Product not found" } }
  { "data": [ ... ], "meta": { "page": 1, "page_size": 20, "total_items": 134, "total_pages": 7 } }
  ```

## Testing conventions

Testify (`assert`/`require`, `suite` where it cuts boilerplate). Per AGENTS.md:
- **Repository tests** run against a real (test) PostgreSQL database — not mocked — since the point is verifying the actual SQL/constraints behave.
- **Service tests** mock the repository *interface* (this is the concrete payoff of DIP — see AGENTS.md's SOLID section) and assert business outcomes: can't oversell stock, receiving increases inventory, cancelling an unreceived purchase leaves inventory untouched, unauthorized actions are rejected.
- **Handler tests** use `httptest` + Gin's test mode, asserting status codes and envelope shape, not internal implementation.
- Table-driven tests for anything with more than two meaningful cases. See `templates/service_test.template.go`.

## Reviewing existing code

When asked to review, or touching a file for an unrelated reason, check for and flag:

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
- [ ] An error logged with `log.Error().Err(err)` directly in a service instead of via `logger.LogError` (loses the dev/prod detail switch)
- [ ] Sensitive fields (password, token, etc.) passed into a logged fields map
- [ ] A handler or endpoint change with no matching Swagger annotation update, or annotations updated without re-running `swag init`
- [ ] A repository interface with no corresponding mock used in service tests

Report findings as a short, concrete list — then apply fixes if asked to fix rather than just review.

## Templates

Read the relevant template before writing the corresponding file — they encode the exact interface shapes, error handling, and logging calls used across this codebase:

- `templates/repository.template.go` — repository interface + GORM implementation
- `templates/service.template.go` — service with DI, sentinel errors, transaction + logging example
- `templates/handler.template.go` — Gin handler using the response package + Swagger annotation
- `templates/response.template.go` — uniform success/error/pagination response helpers
- `templates/logger.template.go` — dev/production logging switch with sensitive-field redaction
- `templates/swagger-regen.template.txt` — how and when to regenerate Swagger docs
- `templates/service_test.template.go` — Testify service test with a mocked repository
