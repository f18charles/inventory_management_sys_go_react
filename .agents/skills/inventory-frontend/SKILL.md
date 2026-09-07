---
name: inventory-frontend
description: >-
  Project-specific conventions for the Inventory Management System's React frontend
  (Vite, TanStack Router, TanStack Query, Zustand, axios, Tailwind, plain
  JavaScript/JSX, no TypeScript). Use this skill whenever writing, scaffolding,
  reviewing, or fixing ANY frontend code in this repo's frontend directory — new
  pages, routes, components, Zustand stores, API hooks, forms, or layouts. Also use
  it when reviewing existing frontend code for convention drift (wrong folder, mixed
  state management, ad-hoc fetch or axios calls outside the shared client,
  inconsistent handling of the backend's response envelope), even if the user didn't
  explicitly ask for a review. Trigger on requests like "add a products page,"
  "create a store for X," "add a hook to fetch Y," "build the sales form," or "does
  this component look right." This skill assumes the backend conventions defined in
  AGENTS.md (Gin, versioned REST API under /api/v1, and a data/error JSON envelope
  on every response) are already in place on the backend being called.
---

# Inventory Frontend Skill

Project-specific frontend conventions for the Inventory Management System. This is scoped to `frontend/` in this repo only — it is not a general React skill.

## When generating new code

Follow this decision path for any new frontend feature:

1. **Does it need a route?** → new file under `src/routes/` (TanStack Router file-based routing). See `templates/route.template.jsx`.
2. **Does it need server data?** → a query/mutation hook under `src/api/` using TanStack Query, calling the shared API client. See `templates/query-hook.template.js` and `templates/api-client.template.js`.
3. **Does it need shared client-side state** (not server data — auth session, UI toggles, multi-step form state)? → a Zustand store under `src/stores/`. See `templates/zustand-store.template.js`. **Do not** put server data (products, sales, inventory) in Zustand — that's TanStack Query's job. Zustand is for client-only state.
4. **Is it a reusable piece of UI** with no route of its own? → `src/components/`.
5. **Is it a full route-level view** composed of components + hooks? → `src/pages/`, imported by the matching file in `src/routes/`.
6. **Is it a shared shell** (e.g. the authenticated app chrome, the auth/login shell)? → `src/layouts/`.

Always check `src/api/`, `src/stores/`, and `src/components/` for something that already does what's needed before creating a new file — don't duplicate an existing hook or store slice.

## Folder structure (authoritative)

```
frontend/src/
├── components/   # reusable, route-agnostic UI
├── pages/        # route-level views (composition only, minimal logic)
├── layouts/      # shared shells (AppLayout, AuthLayout)
├── hooks/        # non-API custom hooks (e.g. useDebounce, useMediaQuery)
├── stores/       # Zustand stores — client state only
├── routes/       # TanStack Router route files (file-based)
├── api/          # API client + TanStack Query hooks (server state)
├── utils/        # pure helper functions (formatCurrency, formatDate, etc.)
└── types/        # JSDoc typedefs for shared shapes (no TypeScript in this repo)
```

## API layer conventions

- **One shared axios instance** (`src/api/client.js`, see `templates/api-client.template.js`): `axios.create` with `baseURL` = `/api/v1`, a request interceptor that attaches the auth token from the auth store, and a response interceptor that unwraps the backend's `data` envelope and normalizes errors.
- Backend responses are always `{ "data": ... }` on success or `{ "error": { "code": "...", "message": "..." } }` on failure (per AGENTS.md). The response interceptor unwraps `response.data.data` on success, and on failure reads `error.response.data.error` and rejects with a normalized `ApiError` (carrying `code`, `message`, `status`) — so calling code never touches axios's `error.response` shape or the raw envelope directly.
- **Never import `axios` directly in a component or call `axios.get/post` ad hoc.** Every backend call goes through the shared instance, wrapped in a TanStack Query hook.
- One hook file per resource (`src/api/products.js`, `src/api/sales.js`, ...), exporting query hooks (`useProducts`, `useProduct(id)`) and mutation hooks (`useCreateProduct`) — see `templates/query-hook.template.js`.
- Query keys are arrays scoped by resource and params, e.g. `['products', { categoryId }]`, so invalidation after a mutation is precise (`queryClient.invalidateQueries({ queryKey: ['products'] })`), not a blanket refetch-everything.
- Money values from the backend are integers (smallest currency unit, per AGENTS.md). Format for display with a shared `formatCurrency` util in `src/utils/` — never do ad-hoc `/100` math inline in a component.

## State management conventions

- **Server state → TanStack Query. Client state → Zustand. Don't mix them.** If data comes from the backend, it does not belong in a Zustand store, even "for convenience" — it goes stale silently and causes bugs that are hard to trace.
- One store per concern, not one giant global store: `useAuthStore` (session/user/token), `useUiStore` (sidebar open, active modal), feature-specific stores only when a flow genuinely needs cross-component client state (e.g. a multi-step sale-creation wizard before it's submitted).
- Stores expose actions alongside state (`{ user, token, login, logout }`), not raw setters sprinkled through components.
- Don't put derived data in a store — compute it in a selector or in the component.

## Routing conventions

- File-based routes under `src/routes/`, generated into a route tree by the TanStack Router Vite plugin. Route files stay thin: they define the path, loaders/guards, and render the matching `src/pages/` component — no business logic in the route file itself.
- Auth-gated routes use a layout route (e.g. `src/routes/_authenticated.jsx`) that checks `useAuthStore` and redirects to login rather than each route re-implementing the check.
- Route params map directly to the resource ID conventions used in the API (`/products/$productId` → `useProduct(productId)`).

## Component conventions

- Function components, hooks only (no class components).
- Props destructured in the function signature; no prop-drilling more than one level — reach for the relevant store or a query hook instead.
- Loading and error states are handled explicitly for every query-backed component (skeleton/spinner for loading, a real error message for the error case) — never render nothing or crash silently while data is in flight.
- Forms: controlled inputs, validation errors shown inline, submit disabled while a mutation is pending, and the mutation's error surfaced to the user (not just logged to console).
- Tailwind utility classes directly in JSX; no separate CSS files per component. For anything about visual design system, spacing scale, or aesthetic direction, defer to the `frontend-design` skill — this skill covers architecture/data flow, not visual styling.

## Reviewing existing code

When asked to review, or when touching a file for an unrelated reason, check for and flag:

- [ ] `fetch` or `axios` called directly in a component instead of going through `src/api/client.js`
- [ ] Server data (products/sales/inventory/etc.) stored in Zustand instead of TanStack Query
- [ ] A new top-level Zustand store created for something that's really route-local state (should be `useState`/`useReducer` instead)
- [ ] Money formatted or divided inline instead of via the shared `formatCurrency` util
- [ ] A query-backed component with no loading state, no error state, or both
- [ ] Business logic (e.g. deciding whether a sale can be submitted) living in a component instead of being driven by backend validation/errors
- [ ] A file in the wrong folder per the structure above (e.g. a reusable component sitting in `pages/`)
- [ ] Route files containing real logic instead of delegating to `pages/`

Report findings as a short list of concrete fixes, not a lecture — then apply them if asked to fix rather than just review.

## Templates

Read the relevant template before writing the corresponding file type — they encode exact import paths and naming conventions used across this codebase, not just illustrative examples:

- `templates/api-client.template.js` — shared axios instance + interceptors + `ApiError`
- `templates/query-hook.template.js` — TanStack Query hooks for a resource
- `templates/zustand-store.template.js` — Zustand store shape/conventions
- `templates/route.template.jsx` — thin TanStack Router route file
- `templates/page.template.jsx` — route-level page composing hooks + components
