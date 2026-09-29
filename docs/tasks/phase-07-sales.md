# Phase 7: Sales

> Source: `docs/Roadmap.md` → Phase 7
> Status: Not started
> Depends on: Phase 4 (product), Phase 5 (inventory)

**What you'll learn:** server-side price recalculation, validating availability across multiple
line items before a write, atomic decrement of inventory during checkout, payment-method modeling
as a constrained enum, and how to design a refund flow on top of an immutable sale record.

Reference: `design.md` §2.3 (customers/sales/sale_items schema), §3.1 (dependency graph),
§3.2 (customer & sale domains), §4 (RBAC), §6.1/§6.3 (atomic decrements, server-side pricing).

## Deliverables

### Customer domain (`internal/customer`)
- [ ] `repository.go` — CRUD.
- [ ] `service.go` — unique customer email handling; sentinel errors.
- [ ] `handler.go` + `routes.go`.
- [ ] `GET`, `POST`, `PUT`, `DELETE /api/v1/customers` (staff may read/create; `manager`/`admin`
      for full writes — see RBAC table).

### Sale domain (`internal/sale`)
- [ ] `repository.go` — create-with-items, get-with-items (preload `Items` and `Customer`), list,
      `UpdateStatus`.
- [ ] `service.go` — `SaleService`.
- [ ] Checkout validation:
  - [ ] Load each product from the database and recalculate
        `sub_total = quantity * product.unit_price` and `total_amount` server-side.
  - [ ] Ignore any client-supplied prices or totals.
  - [ ] Verify available stock for every line item before writing.
- [ ] Atomic completion transaction:
  - [ ] Decrement inventory for all line items and persist sale + items in one transaction.
  - [ ] Roll back and return `ErrInsufficientInventory` if any item is short.
- [ ] Payment method handling: `cash`, `card`, `mobile_money`, `bank_transfer`.
- [ ] `handler.go` — request structs for order + nested items with `binding` tags.
- [ ] `routes.go` — `GET`, `POST /api/v1/sales`; `POST /api/v1/sales/:id/refund`.
- [ ] Role checks per RBAC table (staff may read/create sales; refund restricted to
      `manager`/`admin`).

### Refund / cancellation workflow
- [ ] **Design the rules before coding** (states, which statuses are refundable, whether stock is
      restored, whether partial refunds are in scope for v1) — write them down first.
- [ ] Implement cancellation of a `pending` sale (no stock change).
- [ ] Implement refund of a `completed` sale: restore inventory for line items atomically and move
      status to `refunded`; guard against double-refund (`ErrInvalidStatusTransition`).
- [ ] Record who performed the refund and when (audit logging).

### Wiring & tests
- [ ] Wire customer and sale domains into `cmd/server/main.go`; inject `inventory.Service` and
      `product` access as needed (interfaces, no cycles).
- [ ] Unit-test price recalculation and rejection of client-tampered prices.
- [ ] Test that a sale with insufficient stock is rejected and nothing is written.
- [ ] Test successful checkout decrements stock by the exact quantities.
- [ ] Test refund restores stock exactly once and rejects a second refund.
- [ ] Regenerate Swagger docs (`make swagger`).

## Notes

- Sale + items + inventory decrements must be one atomic transaction.
- Refund rules are business decisions — capture them in `design.md` or the PR description before
  implementing so the behavior is unambiguous.
- Money is `int64` minor units; totals always come from the database.

## Definition of Done

- [ ] Code compiles cleanly; dependency graph stays acyclic.
- [ ] Checkout and refund run in transactions with rollback.
- [ ] All totals recalculated server-side; no float money.
- [ ] Layer boundaries respected.
- [ ] `response` envelopes used; failures logged without secrets.
- [ ] Role checks applied; refund restricted.
- [ ] Refund/cancellation rules documented.
- [ ] Tests pass (`make test`).
- [ ] Swaggo annotations present and `make swagger` re-run.
