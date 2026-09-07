# AGENTS.md

Instructions for AI coding agents working in this repository.

## Project Overview

Full-stack inventory management system for small and medium-sized businesses. It manages products, categories, suppliers, customers, stock levels, purchases, sales, users, and inventory activity through a Go backend and a React frontend.

This is a production-oriented portfolio project **and** a learning project. When implementing a feature, explain the architectural reasoning behind non-trivial decisions if asked — the goal is understanding, not just a compiling result. See `ROADMAP.md` for the phased build-out plan.

## Tech Stack

**Backend:** Go, Gin, GORM, PostgreSQL, `golang-migrate`, REST API, JWT authentication, bcrypt/Argon2 password hashing, Testify for tests, Swagger/OpenAPI for API documentation, `zerolog` for structured logging.

**Frontend:** React, Vite, TanStack Router, TanStack Query, Zustand, axios, Tailwind CSS.

**Infrastructure:** Docker, Docker Compose, PostgreSQL, GitHub Actions, environment-based configuration.

Gin is the chosen web framework — do not fall back to raw `net/http` for routing, and do not swap in a different framework (Echo, Fiber, Chi, etc.) without an explicit decision to change direction. Keep Gin usage disciplined: it lives in the handler layer only (routing, binding, middleware, responses). Business logic and persistence never move into a handler just because Gin makes it easy to inline them there.

TanStack Router is the chosen frontend router and Zustand the chosen global state library — do not introduce Redux, Context-based global state, or a different router without a deliberate decision to change direction.

Before adding any dependency, ask: Do we need it? Does the standard library or an already-adopted tool solve this? Is it maintained? Does it add real complexity? Does it fit this architecture?

## Architecture

The backend is strictly layered. Do not bypass layers without a clear reason.

```
Gin Handler → Service → Repository → GORM → PostgreSQL
```

- **Handler** — HTTP concerns only: routing, request binding/validation (via Gin), extracting auth info, calling services, setting status codes, shaping the JSON response. No database queries, no business logic. Request-level logging happens in middleware, not repeated per handler (see Logging).
- **Service** — Business logic and workflow coordination: can this sale be completed, is there enough inventory, can this purchase be received, is the user authorized, should a transaction be started. Services must not know about HTTP status codes, Gin's `*gin.Context`, or response objects — return domain errors (e.g. `ErrInsufficientInventory`), not `http.StatusBadRequest`.
- **Repository** — Persistence only: execute queries, handle persistence errors, return domain/model data. No HTTP logic, no JWT checks, no business decisions.
- **GORM / PostgreSQL** — Data mapping and storage. GORM does not replace the need to understand SQL.

## Domain Models

```
BaseModel (composed into every model)
User, Supplier, Customer, Category, Product, Inventory,
Purchase, PurchaseItem, Sale, SaleItem
```

`BaseModel` provides `ID` (UUID), `CreatedAt`, `UpdatedAt`. UUID generation happens via a GORM `BeforeCreate` hook — don't duplicate ID-generation logic per model without reason.

Relationships:
```
Category → Products → Inventory
Supplier → Purchases → PurchaseItems → Products
Customer → Sales → SaleItems → Products
User → Purchases, Sales
```

Keep models focused on data representation. Do not put large business workflows in models or GORM hooks (hooks are for small, predictable, persistence-adjacent behavior only — e.g. generating an ID, not "create sale → check inventory → adjust stock → notify").

## Database

- PostgreSQL is the system of record. Enforce integrity at the database level: primary keys, foreign keys, unique constraints, not-null constraints, check constraints, indexes. Do not rely solely on application-level validation.
- Prefer native types where useful: `UUID`, `BOOLEAN`, `TIMESTAMPTZ`, `NUMERIC`.
- Use `golang-migrate` with versioned, numbered migrations (`000001_create_users.up.sql` / `.down.sql`). **Never** use GORM `AutoMigrate` for schema management, and **never** edit a migration that has already run against a shared database — add a new one instead.
- Think about deletion behavior explicitly; don't reach for cascading deletes by default. Historical sales/purchases must stay meaningful even if a referenced product is later deactivated.
- Use GORM preloading intentionally (`Preload("Items")`), not by default — consider query cost, response size, N+1 risk, and whether the endpoint actually needs the relation.

## Business Rules

**Inventory**
- Stock must never go negative.
- A sale can only complete if sufficient stock exists.
- Receiving a purchase increases inventory; cancelling an *unreceived* purchase must not touch inventory.
- Refund-related inventory behavior must be explicitly designed before it's implemented.
- Any inventory-changing operation runs inside a transaction.

**Purchases** — statuses: `pending`, `received`, `cancelled`. Workflow: create → add items → pending → receive → increase inventory → received.

**Sales** — statuses: `pending`, `completed`, `cancelled`, `refunded`. Payment methods: `cash`, `card`, `mobile_money`, `bank_transfer`. Workflow: create → add items → check inventory → confirm → decrease inventory → completed.

**Transactions** — required whenever a business operation touches multiple related records that must succeed or fail together (e.g. sale + sale items + inventory decrement; purchase receipt + inventory increment). No partially completed sales or purchases — roll back on any failure.

**Money** — never use floating-point for currency. Store integers in the smallest currency unit (e.g. `int64` cents; KSh 100.50 → `10050`). Be consistent throughout the backend. Never trust client-provided totals — always recompute `quantity × unit price` on the backend.

## Authentication & Authorization

- Roles: `Admin`, `Manager`, `Staff` (see README for the full permission matrix by resource).
- Passwords: hash with bcrypt or Argon2, never store plaintext, never return `PassHash`/password hashes in any API response.
- JWT-based authentication with token validation middleware; role-based authorization enforced in the backend.
- Never rely on the frontend to enforce permissions — React can provide UX-level checks only; the backend is authoritative.

## API Design

- REST conventions under `/api/v1/...` (e.g. `GET/POST /api/v1/products`, `GET/PATCH/DELETE /api/v1/products/{id}`), implemented as Gin routes/route groups.
- Use correct status codes: `200`, `201`, `204`, `400`, `401`, `403`, `404`, `409`, `422`, `500`.
- Consistent JSON envelopes, built via `internal/utils/response` — handlers never construct `gin.H{...}` directly:
  - Success: `{ "data": { ... } }`
  - Paginated: `{ "data": [ ... ], "meta": { "page": 1, "page_size": 20, "total_items": 134, "total_pages": 7 } }`
  - Error: `{ "error": { "code": "...", "message": "..." } }`
- Never leak raw SQL errors, DB connection details, stack traces, or internal implementation details to clients — this is independent of the logging dev/prod switch, which only affects what gets *logged server-side*, never what's returned to a caller.

## API Documentation

The API is documented using Swagger/OpenAPI.

- Annotate handlers with `swaggo/swag`-style comments so docs are generated from the code rather than maintained by hand in a separate file.
- Every endpoint documents: method, path, request body/params, possible response codes, and response shape (including the standard `data`/`error`/`meta` envelopes).
- **Generation is a workflow step, not automatic** — `swag init` only rebuilds `docs/` when explicitly run; editing an annotation comment does nothing on its own. Wire it into a Makefile target (`make swagger`) and a CI check that fails the build if generated docs are stale against the committed ones. Stale API docs are treated as a bug.
- Serve the generated docs (e.g. `/swagger/index.html`) in development; decide deliberately whether they're exposed in production.

## Validation

Split responsibilities across layers — don't duplicate complex business validation everywhere:
- **Handler:** is the JSON well-formed, are required fields present (Gin binding/validator tags are the natural place for this).
- **Service:** does this make business sense (stock available, user permitted, purchase receivable, etc.).
- **Database:** can this data legally exist (constraints).

## Error Handling

Define reusable sentinel errors, e.g.:
```go
var (
    ErrNotFound              = errors.New("resource not found")
    ErrUnauthorized          = errors.New("unauthorized")
    ErrForbidden             = errors.New("forbidden")
    ErrInsufficientInventory = errors.New("insufficient inventory")
)
```
Handlers translate these into HTTP responses. Never expose SQL syntax errors, DB internals, or stack traces to API clients.

## Logging

Every request and every business-significant event must produce a log line that's readable in the terminal during development and parseable in production. Silence is the enemy here — if a transaction fails and nothing gets logged, the bug is invisible until a customer notices bad stock or a missing sale.

**Library:** `zerolog`, initialized once at startup via `internal/utils/logger`, driven by a single `APP_ENV` switch:
- **Development** (`APP_ENV != production`): colorized `zerolog.ConsoleWriter`, `debug` level, and failure logs (via the shared `logger.LogError` helper) include the full underlying error text — use this when actively debugging.
- **Production** (`APP_ENV == production`): plain JSON lines (so logs can be shipped to a real aggregator later), `info` level, and `logger.LogError` logs a safe summary with `detail_suppressed: true` instead of the raw error text — so an internal error message (which may contain SQL fragments, file paths, or request-derived values) never sits in production logs by default. To see the exact issue, restart with `APP_ENV=development`; that's the intended way to get full detail, not a code change.
- Known-sensitive keys (`password`, `pass_hash`, `token`, `authorization`, `jwt_secret`, `secret`) are stripped from structured log fields in **both** modes, unconditionally.
- Don't build a custom logger from scratch and don't use `log.Println`/`fmt.Println` for anything beyond a throwaway debug print — those don't carry levels, fields, or timestamps and get stripped out before commit.

**Where logging happens:**
- **Middleware** — a request-logging middleware logs every request once it completes: method, path, status code, latency, and (once auth exists) the authenticated user ID. This is the one place logging is automatic and mandatory for every endpoint, not opt-in per handler.
- **Handler** — logs only when something handler-specific goes wrong (bad request body, auth failure) at `warn` level; don't re-log what the middleware already captured.
- **Service** — this is where most *meaningful* logs belong, since it's where business decisions happen. Log the start and outcome of significant workflows — sale created, purchase received, inventory adjusted, insufficient stock rejected — with enough fields (IDs, quantities, amounts) to reconstruct what happened without a debugger.
- **Repository** — logs unexpected persistence errors (`error` level) with the failing operation and relevant IDs, but does not log routine "not found" cases as errors — those are expected outcomes the service/handler translate normally.

**Transactions specifically:** log at the start of a transaction (`debug`, or `info` for high-value operations), and always log the outcome — success with the resulting IDs/amounts, or failure with the error and enough context to know which step failed and that a rollback occurred. A transaction that fails silently is the single worst logging gap in this system — treat it as a bug if it happens.

**Levels:**
- `debug` — verbose, dev-only detail (SQL timing notes, step-by-step trace through a workflow).
- `info` — normal significant events (request completed, sale completed, purchase received, user logged in).
- `warn` — recoverable/expected failure conditions (validation failure, insufficient stock, unauthorized attempt).
- `error` — unexpected failures that need attention (DB errors, transaction rollbacks, panics recovered by middleware).

**Rules:**
- Never log passwords, password hashes, JWTs, or full request bodies containing sensitive data — this is enforced by `internal/utils/logger`'s field redaction, not left to per-call-site discipline.
- Failures go through `logger.LogError(err, safeSummary, fields)`, not a raw `log.Error().Err(err).Msg(...)` — that's what gives the dev/prod detail switch above. A raw zerolog error call in a service is a review flag.
- Use structured fields (`log.Info().Str("sale_id", id).Msg("sale completed")`), not string-concatenated messages — this is what makes logs greppable and parseable later.
- Include a request ID (generated in middleware, passed through context) on every log line tied to a request, so a single transaction's logs can be traced end-to-end even under concurrent load.
- Don't log and then also return the same error up the stack without adding context — log once at the boundary that actually handles it (usually the repository for the raw error, or the handler for the final HTTP outcome), and wrap/annotate the error as it travels rather than logging it at every layer it passes through.

## Go Style & SOLID Principles

Idiomatic Go, standard `if err != nil { return err }` error handling, focused functions. SOLID is applied pragmatically here, in a way that fits Go and the existing layered architecture — not as a rule to satisfy for its own sake, and not a license to add abstraction beyond what's below.

- **Single Responsibility** — largely what the Handler → Service → Repository split already enforces. A handler parses/responds, a service decides, a repository persists. If a function is doing two of those jobs, split it.
- **Open/Closed** — favor adding new behavior (new service methods, new repository methods, new handlers) over modifying shared logic in ways that risk existing behavior. Doesn't mean over-engineering for hypothetical extension points that don't exist yet.
- **Liskov Substitution** — if a repository or service interface has multiple implementations (e.g. a real repository and a test/mock repository), the mock must honor the same contract and error behavior as the real one — no surprising shortcuts.
- **Interface Segregation** — keep interfaces small and specific to what a consumer actually needs (e.g. a service depending on a `ProductReader` rather than a full `ProductRepository` when it only reads). Don't create one large interface "for completeness."
- **Dependency Inversion** — services depend on repository *interfaces*, not concrete GORM implementations, so persistence can be swapped or mocked in tests. This is the main place an interface earns its keep in this codebase — don't create interfaces solely because "interfaces are good"; testing/substitution/architectural boundaries are the concrete reasons that justify one.

In short: SOLID here mainly justifies *why* repository interfaces exist for testability (DIP) and why the layers stay separated (SRP) — it's not a mandate to add abstraction layers beyond that.

## Configuration & Secrets

- Config comes from environment variables / uncommitted config files.
- Never hardcode DB passwords, JWT secrets, API keys, or production credentials.
- `.env` is git-ignored; provide `.env.example` when useful.
- In Docker, remember `localhost` refers to the current container — reach PostgreSQL via its Compose service name.

## Testing

New functionality needs tests, written with Testify (`assert`/`require`, and `suite` where it reduces boilerplate). Prioritize business rules, repository behavior, transactions, auth, inventory changes, sales, and purchases. Test behavior, not implementation details (don't just assert that function A calls function B).

- **Model tests:** UUID generation, hooks, defaults.
- **Repository tests:** CRUD, relationships, query filtering.
- **Service tests:** e.g. can't oversell stock, receiving increases inventory, cancelling an unreceived purchase doesn't touch inventory, unauthorized actions are blocked.
- **Handler tests:** status codes, request parsing/binding, response bodies, auth behavior, invalid input (`httptest` + Gin's test mode).
- **Integration tests:** multiple layers together against a real test PostgreSQL database.

Goal is confidence in important behavior, not raw test count.

## Git Workflow

- Feature branches: `feat/product-management`, `feat/sales`, `fix/inventory-update`, `refactor/repository-layer`, etc.
- Conventional, descriptive commits: `feat: add product repository`, `fix: prevent negative stock`, `test: add inventory service tests`, `refactor: separate sale business logic`.
- Never commit `.env`, credentials/secrets, build artifacts, IDE files, or temp files.

## Feature Implementation Order

```
Requirement → Domain/model consideration → Migration
   → Repository → Repository tests
   → Service → Service tests
   → Handler → Handler tests → Swagger annotations
   → API integration → Frontend (if applicable)
```

For multi-record database changes: Service → Transaction → Repository operations → Commit/Rollback.

## Before Writing Code

Inspect the existing implementation first — current models, migrations, DB config, repository/service patterns, existing error definitions, test utilities, and API conventions. Don't recreate functionality that already exists because the repo wasn't checked.

## Before Finishing a Task — Checklist

- [ ] Code compiles
- [ ] Migration runs and rolls back cleanly
- [ ] Tests pass (run the relevant suite before declaring done)
- [ ] Errors are handled and don't leak internals
- [ ] Relationships/foreign keys are correct
- [ ] Transactions are used where required
- [ ] API responses built via `internal/utils/response`, following the standard envelope (including `meta` for paginated lists)
- [ ] New/changed endpoints have Swagger annotations, and `swag init` was re-run so `docs/` isn't stale
- [ ] Failures logged via `logger.LogError` (not a raw zerolog call), with useful fields and no silent failures
- [ ] No secrets or sensitive data end up in logs
- [ ] No unnecessary coupling or abstraction was introduced

## Changes to Existing Architecture

Understand *why* a boundary exists before crossing it. Specifically avoid:
- Moving DB queries into handlers for convenience.
- Putting business logic in repositories.
- Putting HTTP logic in services.
- Adding global state without a clear reason.
- Adding interfaces/abstractions without a concrete, present problem.

Prefer the simplest architecture that keeps responsibilities clearly separated.