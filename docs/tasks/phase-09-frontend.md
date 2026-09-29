# Phase 9: Frontend

> Source: `docs/Roadmap.md` → Phase 9
> Status: Not started
> Depends on: Phases 4–8 (working, documented API)

**What you'll learn:** how Vite structures a React app, how TanStack Router defines a typed route
tree and layouts, how TanStack Query handles server state (keys, caching, mutations, invalidation),
how Zustand holds session/UI state, and how a shared Axios client attaches the JWT and unwraps the
response envelope.

Reference: `design.md` §1 (frontend stack), §5 (file structure), `AGENTS.md` §2 (frontend
conventions). Plain JavaScript/JSX — no TypeScript.

## Backend prerequisite: Reports domain

The Reports pages need `internal/report` endpoints (`design.md` §3.2/§4). The Roadmap does not list
them under a backend phase, so build them here before the Reports screens.

- [ ] `internal/report/repository.go` — aggregate queries (totals, low-stock count, top products).
- [ ] `internal/report/service.go` — shape summaries and date-range filtering.
- [ ] `internal/report/handler.go` + `routes.go` — `GET /api/v1/reports/summary`,
      `GET /api/v1/reports/sales`; role-restricted to `manager`/`admin`.
- [ ] Wire into `cmd/server/main.go`; annotate and regenerate Swagger.

## Setup

- [ ] Scaffold the Vite React app in `frontend/` (plain JS/JSX).
- [ ] Install and configure Tailwind CSS.
- [ ] Install TanStack Router, TanStack Query, Zustand, and Axios.
- [ ] Configure `vite.config.js` (dev proxy or API base URL via env, e.g. `VITE_API_URL`).
- [ ] Add `npm run dev`, `build`, and `lint` scripts; wire the Makefile targets
      (`make frontend-dev`, `make frontend-build`, `make frontend-lint`).

## Core plumbing

- [ ] `src/api/client.js` — Axios instance with a request interceptor that attaches the JWT and a
      response interceptor that unwraps the `data`/`error` envelope and handles `401` (logout).
- [ ] `src/stores/useAuthStore.js` — Zustand store for token, user, role, and login/logout actions.
- [ ] `src/routes` — TanStack Router route tree with `AuthLayout` (public) and `AppLayout`
      (protected) and role-based route guards.
- [ ] `src/layouts` — `AppLayout`, `AuthLayout`, `Sidebar`, `Navbar`.
- [ ] `src/hooks` — TanStack Query hooks per domain (`useProducts`, `useCategories`, `useInventory`,
      `useSuppliers`, `useCustomers`, `usePurchases`, `useSales`, `useReports`, `useUsers`).
- [ ] `src/components` — shared table, modal, form field, button, and badge components.

## Screens

- [ ] Authentication flow — login screen wired to `POST /api/v1/auth/login` and session bootstrap
      via `GET /api/v1/auth/me`.
- [ ] Dashboard — KPI summary cards and sales trend charts from report endpoints.
- [ ] Products — paginated table with search, category filter, and create/edit modal.
- [ ] Categories — list and create/edit modal.
- [ ] Inventory — monitoring table with low-stock badges and a quick adjustment dialog.
- [ ] Suppliers — CRUD list and form.
- [ ] Customers — CRUD list and form.
- [ ] Purchases — order builder, line items, and one-click receiving.
- [ ] Sales / POS — checkout with instant stock validation and payment method selection.
- [ ] Reports — summary metrics, top products, and date-range filtering.
- [ ] User management — list, create, and role updates (admin only).

## Cross-cutting

- [ ] Loading, empty, and error states on every data view.
- [ ] Mutations invalidate the correct query keys so lists refresh after writes.
- [ ] Role-aware UI: hide/disable actions the backend would reject (never as the only guard —
      the backend remains authoritative).
- [ ] Consistent formatting for money (convert minor units to display currency) and dates.

## Definition of Done

- [ ] `make frontend-build` succeeds and `make frontend-lint` passes.
- [ ] All server state goes through TanStack Query; no ad-hoc `fetch`/`axios` calls in components.
- [ ] All client state that is global lives in Zustand; local UI state stays in components.
- [ ] Hooks live in `src/hooks`, stores in `src/stores`, components in `src/components`.
- [ ] Protected routes are guarded and the auth token is attached to API calls.
- [ ] No TypeScript files introduced.
