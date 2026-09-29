# Phase 6: Purchases

> Source: `docs/Roadmap.md` → Phase 6
> Status: Not started
> Depends on: Phase 4 (product), Phase 5 (inventory)

**What you'll learn:** how to persist a parent order with child line items in one transaction, how
GORM handles associations and bulk inserts, how to make a state transition idempotent, and how a
domain service calls another domain's interface (`inventory`) without creating import cycles.

Reference: `design.md` §2.2 (relationships), §2.3 (suppliers/purchases/purchase_items schema),
§3.1 (acyclic dependencies), §3.2 (supplier & purchase domains), §4 (RBAC), §6.2 (receiving
idempotency).

## Deliverables

### Supplier domain (`internal/supplier`)
- [ ] `repository.go` — CRUD.
- [ ] `service.go` — unique supplier email handling; sentinel errors.
- [ ] `handler.go` + `routes.go`.
- [ ] `GET`, `POST`, `PUT`, `DELETE /api/v1/suppliers` (`GET` readable by all roles; writes by
      `manager`/`admin`).

### Purchase domain (`internal/purchase`)
- [ ] `repository.go` — `PurchaseRepository` with create-with-items, get-with-items (preload
      `Items` and `Supplier`), list, and `UpdateStatus`.
- [ ] `service.go` — `PurchaseService`.
- [ ] Create a purchase with line items: validate supplier and each product exist; compute
      `subtotal = quantity * unit_cost` and `total_amount` **server-side** from the database;
      ignore client-supplied totals.
- [ ] Persist purchase + items atomically in one transaction.
- [ ] `handler.go` — request structs for order + nested items with `binding` tags.
- [ ] `routes.go` — `GET`, `POST /api/v1/purchases`; `POST /api/v1/purchases/:id/receive`.
- [ ] State transitions: `pending` → `received` or `cancelled`.
- [ ] Receiving workflow:
  - [ ] Guard that current status is `pending`; reject if already `received`/`cancelled`
        (`ErrInvalidStatusTransition`).
  - [ ] Increment inventory for every line item atomically via `inventory.Service`.
  - [ ] Update purchase status to `received` in the same transaction; roll back on any failure.
  - [ ] Cancelling an unreceived order must not touch stock.
- [ ] Role checks: purchases readable/writable by `manager`/`admin`; receive by `manager`/`admin`.

### Wiring & tests
- [ ] Wire supplier and purchase domains into `cmd/server/main.go`; pass `inventory.Service` into
      `NewPurchaseService` (depend on the interface, not the concrete type).
- [ ] Unit-test total recalculation and rejection of tampered totals.
- [ ] Test receiving increments stock exactly once and that a second receive is rejected.
- [ ] Test cancellation leaves stock unchanged.
- [ ] Regenerate Swagger docs (`make swagger`).

## Notes

- Depend only on `inventory`/`product`/`supplier`; never create a reverse dependency.
- Receiving must be idempotent — status check plus transaction guards prevent double-counting.
- All money is `int64` minor units.

## Definition of Done

- [ ] Code compiles cleanly; dependency graph stays acyclic.
- [ ] Purchase + items + inventory updates run in one transaction with rollback.
- [ ] Totals and subtotals recalculated server-side.
- [ ] Layer boundaries respected.
- [ ] `response` envelopes used; failures logged without secrets.
- [ ] Role checks applied on write/receive routes.
- [ ] Tests pass (`make test`).
- [ ] Swaggo annotations present and `make swagger` re-run.
