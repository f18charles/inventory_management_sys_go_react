# Task Lists

Per-phase task breakdown for everything left in `docs/Roadmap.md`. Each file is a self-contained
checklist you can work through top to bottom. The backend files are written to build fluency in
**Go / Gin / GORM**; the frontend files are written to build fluency in **React** (Vite, TanStack
Router, TanStack Query, Zustand).

Read `AGENTS.md` (ground rules), `design.md` (architecture and schema), and `PRD.md` (scope)
alongside these lists. Do not start a phase until the previous one compiles, migrates, and passes
its tests.

## Remaining phases

| Phase | Task List | Focus | Status |
|---|---|---|---|
| 4 | [Product Management](./phase-04-product-management.md) | Go/Gin/GORM — category & product domains | Not started |
| 5 | [Inventory](./phase-05-inventory.md) | Go/Gin/GORM — stock service & transactions | Not started |
| 6 | [Purchases](./phase-06-purchases.md) | Go/Gin/GORM — suppliers & receiving | Not started |
| 7 | [Sales](./phase-07-sales.md) | Go/Gin/GORM — customers & checkout | Not started |
| 8 | [API Documentation](./phase-08-api-documentation.md) | Go/Swaggo — OpenAPI annotations | Not started |
| 9 | [Frontend](./phase-09-frontend.md) | React — Router, Query, Zustand | Not started |
| 10 | [Production Readiness](./phase-10-production-readiness.md) | Go + tooling — tests, CI, Docker | Not started |
| — | [Technology Learning Track](./learning-track.md) | Redis, Kafka, Sentry, PostHog, BI, AWS | Not started |

## Completed phases

| Phase | Focus | Status |
|---|---|---|
| 1 | Database foundation & migrations | Done (seed data pending) |
| 2 | Backend foundation & middleware | Done |
| 3 | Authentication & user management | Done |

## How to use these lists

1. Pick the next incomplete phase.
2. Create a branch: `<type>/<short-slug>` (e.g. `feat/phase-4-product-management`).
3. Work one checkbox at a time; commit logical chunks with clear messages.
4. Do not tick a box until the layer compiles and its tests pass.
5. Finish by running the phase's **Definition of Done** checklist before moving on.
