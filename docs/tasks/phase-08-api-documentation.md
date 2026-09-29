# Phase 8: API Documentation

> Source: `docs/Roadmap.md` → Phase 8
> Status: Not started
> Depends on: Phases 4–7 (endpoints exist to document)

**What you'll learn:** how Swaggo turns Go comments into an OpenAPI spec, how to annotate Gin
handlers with `@Summary`/`@Param`/`@Success`/`@Failure`/`@Router`, how to model security with a
Bearer token definition, and how to keep generated docs in sync in CI.

Reference: `design.md` §1 (Swagger tooling), `AGENTS.md` §2 (API documentation convention),
`PRD.md` §7 (100% route coverage).

## Deliverables

### Annotations
- [ ] Add package-level API metadata (title, version, base path `/api/v1`, description).
- [ ] Add a `Bearer` JWT security definition and reference it on protected routes.
- [ ] Annotate every handler across all domains with:
  - [ ] `@Summary`, `@Description`, `@Tags`
  - [ ] `@Accept` / `@Produce`
  - [ ] `@Param` for path, query, and body (referencing request DTOs)
  - [ ] `@Success` / `@Failure` with the standard response envelope types
  - [ ] `@Router` with the HTTP method and path
- [ ] Ensure request/response DTOs have example tags where helpful.
- [ ] Confirm the standard envelope (`data` / `error` / `meta`) is reflected in the schemas.

### Generation & serving
- [ ] Run `make swagger` (`swag init -g cmd/server/main.go -o docs`) and commit the generated
      `backend/docs` output.
- [ ] Serve Swagger UI at `/swagger/index.html` **in development mode only**, gated by the
      environment config (do not expose docs in production unless explicitly enabled).
- [ ] Verify the UI loads and every `/api/v1` route is present.

### Verification
- [ ] Add a CI check that regenerates docs and fails if they drift from what's committed.
- [ ] Spot-check that annotations map to the real status codes returned by each handler.

## Definition of Done

- [ ] `make swagger` runs without errors and generated files are committed.
- [ ] Every `/api/v1` route appears in the spec with correct method and path.
- [ ] Protected routes show the Bearer security requirement.
- [ ] Swagger UI reachable in development.
- [ ] Docs regenerated as part of the checklist whenever handler annotations change.
