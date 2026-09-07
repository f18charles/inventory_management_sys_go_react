# Roadmap

This file tracks what gets built, in what order. It's split into **Core Build** (the application itself) and **Technology Learning Track** (additional tools layered on once the core is solid, picked up one at a time). Items marked **(manual setup required)** need an account, subscription, or cloud console step from you before an agent can wire up the corresponding code — see `README.md` → "External Services & Manual Setup" for details on those.

Rule for every phase: don't start the next one until the current phase compiles, migrates, and has passing tests where applicable.

---

## Core Build

### Phase 1: Database
- [x] PostgreSQL setup
- [ ] Migration system (`golang-migrate`)
- [ ] Users migration
- [ ] Categories migration
- [ ] Suppliers migration
- [ ] Customers migration
- [ ] Products migration
- [ ] Inventory migration
- [ ] Purchases migration
- [ ] Purchase items migration
- [ ] Sales migration
- [ ] Sale items migration
- [ ] Database indexes
- [ ] Foreign keys
- [ ] Database constraints
- [ ] Seed data

### Phase 2: Backend Foundation
- [ ] Configuration loading (env vars)
- [ ] PostgreSQL connection
- [ ] GORM setup
- [ ] Gin server + router setup
- [ ] Route groups (`/api/v1/...`)
- [ ] Error handling middleware
- [ ] Request validation (Gin binding)
- [ ] Response formatting (`data` / `error` envelope)
- [ ] Structured logger setup (`zerolog`, console output in dev / JSON in prod)
- [ ] Request-logging middleware (method, path, status, latency, request ID)
- [ ] Panic-recovery middleware that logs at `error` level instead of crashing silently
- [ ] Middleware (CORS, etc.)

### Phase 3: Authentication
- [ ] User repository
- [ ] Password hashing (bcrypt/Argon2)
- [ ] Login endpoint
- [ ] JWT authentication
- [ ] Authentication middleware
- [ ] Role-based authorization middleware
- [ ] User management endpoints

### Phase 4: Product Management
- [ ] Category repository / service / handlers
- [ ] Product repository / service / handlers
- [ ] Product search
- [ ] Product filtering by category

### Phase 5: Inventory
- [ ] Inventory repository
- [ ] Stock service
- [ ] Low-stock detection
- [ ] Inventory transaction handling
- [ ] Inventory history

### Phase 6: Purchases
- [ ] Supplier management
- [ ] Purchase repository / service / handlers
- [ ] Purchase items
- [ ] Purchase transactions
- [ ] Receiving workflow

### Phase 7: Sales
- [ ] Customer management
- [ ] Sale repository / service / handlers
- [ ] Sale items
- [ ] Stock validation
- [ ] Sale transactions
- [ ] Refund workflow (design rules explicitly before implementing)

### Phase 8: API Documentation
- [ ] Swagger/OpenAPI annotations on all endpoints
- [ ] Generated spec served at `/swagger/index.html` in development

### Phase 9: Frontend
- [ ] React + Vite setup
- [ ] TanStack Router setup
- [ ] TanStack Query setup (API data fetching/caching)
- [ ] Zustand store(s) for global client state (auth/session, UI state)
- [ ] Authentication flow
- [ ] Dashboard
- [ ] Products
- [ ] Categories
- [ ] Inventory
- [ ] Suppliers
- [ ] Customers
- [ ] Purchases
- [ ] Sales
- [ ] Reports
- [ ] User management

### Phase 10: Production Readiness
- [ ] Automated tests (Testify) across model/repository/service/handler layers
- [ ] Integration tests against a real test database
- [ ] GitHub Actions CI
- [ ] Docker Compose for full local stack
- [ ] Production Docker image
- [ ] Health check endpoint
- [ ] Database backups strategy
- [ ] Deployment
- [ ] Monitoring

---

## Technology Learning Track

Picked up **one at a time**, after Phase 10 is reasonably stable, in roughly this order. Each should land as its own scoped feature branch with a short write-up of why it was added and what it replaced or improved.

### Step A: Redis
- [ ] Add Redis via Docker Compose for local dev
- [ ] Cache layer for read-heavy lookups (e.g. product/category lookups)
- [ ] Rate limiting middleware backed by Redis
- [ ] (Optional) JWT/session blacklist on logout

### Step B: Kafka
- [ ] Add Kafka via Docker Compose for local dev
- [ ] Producer: emit events on sale completion / inventory changes
- [ ] Consumer: build a simple audit-log service off the event stream
- [ ] Document the chosen delivery/consistency guarantees

### Step C: Sentry **(manual setup required)**
- [ ] Backend error tracking integration
- [ ] Frontend error tracking integration
- [ ] Verify errors surface with useful context (no leaked secrets/PII)

### Step D: PostHog **(manual setup required)**
- [ ] Frontend analytics integration
- [ ] Track key product events (login, sale created, purchase received)
- [ ] Confirm no sensitive customer/financial data is sent as event properties

### Step E: Business Intelligence
- [ ] Stand up a BI tool (e.g. Metabase) via Docker for local exploration
- [ ] Connect it (read-only) to the PostgreSQL database
- [ ] Build 2–3 real reports from the Reports feature data (sales trends, low stock, profit)

### Step F: AWS **(manual setup required)**
- [ ] Decide on target services (e.g. RDS for Postgres, S3 for backups/assets, ECS/EC2 for hosting)
- [ ] Infrastructure-as-code or documented manual setup steps
- [ ] Deploy a working environment end-to-end
- [ ] Wire up secrets via AWS-native secrets management rather than committed env files

---

## Notes

- Order within the Technology Learning Track is a suggestion, not a hard rule — reorder if a specific job-relevant tool becomes a priority.
- Don't introduce a Learning Track technology mid-way through a Core Build phase; finish the feature first.
- Every new technology should be wrapped behind the existing architecture (a cache client sits behind a repository/service, not called ad hoc from a handler) — see `AGENTS.md` for the layering rules.