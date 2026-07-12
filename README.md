# production-backend-go

A production-grade blog API in Go — not a tutorial project.

Every architectural decision in this codebase has a documented
reason. The structure is optimized for maintainability over
simplicity, and for operational correctness over convenience.

## Architecture

Vertical slices, not horizontal layers. Each feature owns its
full stack in one place:

```
internal/
├── auth/          # register, login, JWT issuance
│   ├── handler.go    # HTTP layer
│   ├── service.go    # business logic
│   ├── repo.go       # database queries (raw SQL, pgx)
│   ├── routes.go     # route registration
│   └── messages.go   # request/response types
├── post/          # blog post CRUD
│   ├── ...           # same structure as auth
│   └── service_test.go
├── middleware/    # cross-cutting concerns
│   ├── jwt.go        # token validation
│   ├── request_id.go # per-request trace ID
│   ├── recover.go    # panic → 500, never crash
│   └── logger.go     # structured request logging
├── model/         # shared domain types
├── apperr/        # typed application errors
└── platform/      # config, db pool, server wiring
```

When you change the `post` feature, you open one folder.
No cross-cutting file changes for a single feature.

## Operational design

**Two health endpoints:**
- `GET /healthz` — liveness: process is running
- `GET /readyz` — readiness: DB is reachable (returns 503 if not)

Kubernetes uses exactly this distinction. `/readyz` failing
removes the pod from load balancer rotation without killing it.

**`GET /metrics`** — Prometheus scrape endpoint, wired at startup.

**All HTTP timeouts explicitly set** — `ReadHeaderTimeout`,
`ReadTimeout`, `WriteTimeout`, `IdleTimeout`, `ShutdownTimeout`.
Zero timeouts are a DoS vector. None are zero here.

**pgxpool configured for production** — `MaxConns`, `MinConns`,
`MaxConnLifetime`, `MaxConnIdleTime` all explicit. Connection
exhaustion is a defined failure mode, not a surprise.

**Graceful shutdown** — `SIGINT`/`SIGTERM` triggers a shutdown
with a configurable drain timeout. In-flight requests complete.

## Architectural Decision Records

Every non-obvious decision is documented in `docs/adr/`:

| ADR | Decision |
|-----|----------|
| 0001 | Vertical slices over horizontal layers |
| 0002 | Raw SQL with pgx over ORM |
| 0003 | JWT storage strategy |
| 0004 | Migration strategy |
| 0005 | Multi-stage Dockerfile build strategy |

## Runbooks

Operational playbooks in `docs/runbooks/`:
- `migration-failed.md` — what to do when a migration fails
- `database-unavailable.md` — what to do when PostgreSQL is unreachable

## Stack

- Go 1.22
- PostgreSQL via pgx v5 (raw SQL, no ORM)
- chi router
- golang-migrate for schema versioning
- Prometheus metrics
- Docker + Docker Compose (local and prod configs)
- log/slog for structured JSON logging

## Running locally

```bash
cp .env.example .env
make up           # start postgres
make migrate-up   # run migrations
make run          # start server
```

## Quality gate

```bash
make check
# runs: fmt → vet → test → test-race → staticcheck → govulncheck
```

All checks must pass. The pipeline is not decorative.
