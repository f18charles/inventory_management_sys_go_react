# Phase 10: Production Readiness

> Source: `docs/Roadmap.md` → Phase 10
> Status: Not started
> Depends on: Phases 4–9 (feature-complete app)

**What you'll learn:** how to structure Testify unit and integration tests across layers, how to run
tests against a real PostgreSQL instance, how to build CI with GitHub Actions, and how to package
the backend and frontend into a deployable Docker stack with health checks and backups.

Reference: `design.md` §7 (environments, logging, testing), `PRD.md` §7 (quality standards),
`AGENTS.md` §3 (Definition of Done), `WORKFLOW.md` (CI expectations).

## Testing

- [ ] Model tests — validation, UUID generation, enum constraints.
- [ ] Repository tests — against a real test database (`TEST_DATABASE_URL`), including
      transactions and soft-delete behavior.
- [ ] Service tests — business rules, sentinel errors, transaction rollback paths.
- [ ] Handler tests — `httptest` + Gin test mode, asserting the response envelope and status codes.
- [ ] Middleware tests — JWT auth, `RequireRole`, request ID, recovery.
- [ ] Integration tests — full request paths for purchase receiving and sale checkout.
- [ ] Aim for ≥80% coverage across service, repository, and middleware layers
      (`make test-coverage`).

## Continuous integration

- [ ] GitHub Actions workflow that runs `go test -race ./...`, `go vet`/lint, and Swagger drift
      checks.
- [ ] CI step for `npm ci`, `npm run lint`, and `npm run build`.
- [ ] Spin up a PostgreSQL service container for integration tests.

## Containers & local stack

- [ ] Multi-stage production Dockerfile for the backend (build stage + slim runtime).
- [ ] Multi-stage Dockerfile / static build for the frontend.
- [ ] `docker-compose.yml` for the full local stack (PostgreSQL + backend + frontend), with
      migrations applied on startup.
- [ ] `.env.example` documents every required variable; no secrets committed.

## Operations

- [ ] Health check endpoint (e.g. `GET /health` or `/api/v1/health`) reporting DB connectivity.
- [ ] Database backup strategy (scheduled `pg_dump` or managed-provider snapshots) documented.
- [ ] Deployment steps documented and rehearsed end-to-end.
- [ ] Monitoring/logging plan — where errors and metrics are observed in production.

## Definition of Done

- [ ] CI is green on `main`.
- [ ] Coverage target met for the required layers.
- [ ] `docker compose up` brings the stack up from scratch, including migrations.
- [ ] Health check responds and reflects DB status.
- [ ] No secrets or connection strings committed; production errors stay sanitized.
- [ ] Deployment and backup procedures written down.
