# Technology Learning Track

> Source: `docs/Roadmap.md` → Technology Learning Track (Steps A–F)
> Status: Not started
> Depends on: Phase 10 reasonably stable

**What you'll learn:** how to add infrastructure behind existing boundaries without leaking it into
handlers — a cache client behind a repository/service, an event producer behind domain events, and
observability/analytics tools wired so they never expose secrets or PII.

Pick these up **one at a time**, each on its own branch with a short write-up of why it was added
and what it replaced or improved. Don't start one mid-way through a Core Build phase.

## Step A: Redis
- [ ] Add Redis via Docker Compose for local dev.
- [ ] Cache read-heavy lookups (product/category) behind the repository/service layer.
- [ ] Add rate-limiting middleware backed by Redis.
- [ ] (Optional) JWT/session blocklist on logout.
- [ ] Cache invalidation strategy documented (TTL and/or write-through).

## Step B: Kafka
- [ ] Add Kafka via Docker Compose for local dev.
- [ ] Producer: emit events on sale completion and inventory changes.
- [ ] Consumer: simple audit-log service off the event stream.
- [ ] Document delivery/consistency guarantees (at-least-once, ordering, idempotency).

## Step C: Sentry *(manual setup required)*
- [ ] Backend error-tracking integration.
- [ ] Frontend error-tracking integration.
- [ ] Verify errors surface with useful context and no leaked secrets/PII.

## Step D: PostHog *(manual setup required)*
- [ ] Frontend analytics integration.
- [ ] Track key product events (login, sale created, purchase received).
- [ ] Confirm no sensitive customer/financial data is sent as event properties.

## Step E: Business Intelligence
- [ ] Stand up a BI tool (e.g. Metabase) via Docker for local exploration.
- [ ] Connect it read-only to PostgreSQL.
- [ ] Build 2–3 real reports from the Reports feature data (sales trends, low stock, profit).

## Step F: AWS *(manual setup required)*
- [ ] Decide target services (e.g. RDS, S3, ECS/EC2).
- [ ] Infrastructure-as-code or documented manual setup steps.
- [ ] Deploy a working environment end-to-end.
- [ ] Use AWS-native secrets management instead of committed env files.

## Definition of Done

- [ ] Each technology is wrapped behind the existing architecture (no ad-hoc calls from handlers).
- [ ] Local dev bring-up documented in the README.
- [ ] No secrets or PII leaked into logs, events, or analytics properties.
- [ ] A short write-up exists for each step (why, what it replaced, what it improved).
