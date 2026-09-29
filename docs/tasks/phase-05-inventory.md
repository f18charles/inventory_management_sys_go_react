# Phase 5: Inventory

> Source: `docs/Roadmap.md` → Phase 5
> Status: Not started
> Depends on: Phase 4 (product domain)

**What you'll learn:** how to model a one-to-one product↔inventory relation in GORM, how to write
an atomic stock adjustment inside a transaction, how a database `CHECK (quantity >= 0)` constraint
backs up service-level validation, and how to expose low-stock queries.

Reference: `design.md` §2.3 (inventories schema), §3.2 (inventory domain), §4 (RBAC),
§6 (concurrency & atomicity).

## Deliverables

### Repository (`internal/inventory/repository.go`)
- [ ] `InventoryRepository` with `GetByProductID`, `GetByID`, `List`, `Create`.
- [ ] `AdjustQuantity(tx, productID, delta)` — apply a signed change with an atomic SQL update
      (e.g. `UPDATE ... SET quantity = quantity + ? WHERE ...`) so concurrent writes can't race.
- [ ] `ListLowStock` — query products where `quantity <= reorder_level`.
- [ ] Preload the related `Product` (and its `Category`) where useful for display.

### Service (`internal/inventory/service.go`)
- [ ] `InventoryService` exposing `Get`, `List`, `Adjust`, `LowStock`.
- [ ] `Adjust` validates the resulting quantity is non-negative before writing; return
      `ErrInsufficientInventory` when a decrement would go below zero.
- [ ] Wrap the read-check-write in a database transaction; rely on the DB constraint as the final
      safeguard and translate the constraint violation into the sentinel error.
- [ ] Define/reuse sentinel errors in `internal/models/errors.go`.

### Handler & routes
- [ ] `handler.go` — request structs with `binding` tags; responses via `internal/utils/response`.
- [ ] `routes.go` — `GET /api/v1/inventory`, `PATCH /api/v1/inventory/:id/adjust`.
- [ ] `GET /api/v1/inventory` — paginated list with current `quantity` and `reorder_level`.
- [ ] `GET /api/v1/inventory/low-stock` — items at or below reorder level.
- [ ] `PATCH /api/v1/inventory/:id/adjust` — body carries a signed `delta` (or target) plus a
      reason/note; role-restricted to `manager`/`admin`.
- [ ] Log every adjustment through `logger.LogError`/structured logger with the acting user id,
      product id, delta, and a human-readable summary (audit trail).

### Inventory history
- [ ] Decide and document the history source (e.g. derive from purchase/sale items, or add an
      `inventory_adjustments` table via a new migration).
- [ ] If a new table is needed: add matching `.up.sql`/`.down.sql` under
      `backend/internal/database/migrations` and a GORM model in `internal/models`.
- [ ] Expose adjustment history in the list/detail responses.

### Wiring & tests
- [ ] Wire repository/service/handler into `cmd/server/main.go`.
- [ ] Unit-test `Adjust` for increments, valid decrements, and insufficient-stock rejection.
- [ ] Test that concurrent-style double decrement cannot push quantity below zero.
- [ ] Test `ListLowStock` boundary (`quantity == reorder_level` counts as low).
- [ ] Regenerate Swagger docs (`make swagger`).

## Notes

- Never let stock go negative — service check first, DB check constraint as the last line of defense.
- Stock changes should originate from purchase receipt, sale completion, cancellation, or an
  authorized manual adjustment only.

## Definition of Done

- [ ] Code compiles cleanly; no cyclic imports.
- [ ] Multi-record/quantity mutations run inside a transaction.
- [ ] Layer boundaries respected.
- [ ] `response` envelopes used for every response.
- [ ] Failures logged with structured fields; no secrets.
- [ ] Role checks applied (`PATCH` restricted to `manager`/`admin`).
- [ ] Tests pass (`make test`); migrations pair up correctly if added.
- [ ] Swaggo annotations present and `make swagger` re-run.
