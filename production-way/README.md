# go-user-crud

Industry-standard Go REST API — zero framework, pure stdlib routing (Go 1.22+),
clean layered architecture, production-ready from day one.

---

## Architecture

```
cmd/api/              ← main(), config, server bootstrap
internal/
  domain/             ← entities, repository interface, sentinel errors  (no deps)
  service/            ← business logic                                    (depends on domain)
  handler/            ← HTTP encode/decode, routing                       (depends on service interface)
  middleware/         ← request ID, logger, recover, CORS, rate-limit
  repository/
    postgres/         ← pgx/v5 implementation of domain.UserRepository
pkg/
  logger/             ← slog wrapper
  response/           ← JSON helpers, error envelope, pagination
  validator/          ← go-playground/validator singleton
migrations/           ← SQL up/down files
```

**Dependency flow** (strictly one direction):
```
handler → service interface → domain ← repository implements domain
```
No layer imports the one above it. The domain package has zero external dependencies.

---

## Prerequisites

| Tool | Install |
|---|---|
| Go 1.22+ | https://go.dev/dl |
| Docker & Compose | https://docs.docker.com/get-docker |
| `migrate` CLI | `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| `golangci-lint` | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| `air` (optional) | `go install github.com/air-verse/air@latest` |

---

## Quick Start

### 1. Clone & configure

```bash
git clone https://github.com/yourorg/go-user-crud
cd go-user-crud

cp .env.example .env
# Edit .env — at minimum set DATABASE_URL
```

### 2. Start Postgres

```bash
# Option A — Docker Compose (recommended)
make docker/up

# Option B — local Postgres
createdb appdb
```

### 3. Run migrations

```bash
source .env
make migrate/up
```

### 4. Run the server

```bash
# Standard run
make run

# Live reload (requires air)
make watch
```

The API is now available at **http://localhost:8080**.

---

## Project Setup From Scratch (all commands)

```bash
# 1. Create the module
mkdir go-user-crud && cd go-user-crud
go mod init github.com/yourorg/go-user-crud

# 2. Create the directory tree
mkdir -p cmd/api \
         internal/{domain,handler,middleware,service} \
         internal/repository/postgres \
         pkg/{logger,response,validator} \
         migrations \
         .github/workflows

# 3. Install dependencies
go get github.com/jackc/pgx/v5@latest
go get github.com/go-playground/validator/v10@latest
go get github.com/google/uuid@latest

# 4. Install dev tools (into your $GOPATH/bin — not go.mod)
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest
go install golang.org/x/tools/cmd/goimports@latest

# 5. Tidy
go mod tidy
```

---

## API Reference

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/api/v1/users` | List users (paginated) |
| `POST` | `/api/v1/users` | Create user |
| `GET` | `/api/v1/users/{id}` | Get user by ID |
| `PUT` | `/api/v1/users/{id}` | Update user (partial) |
| `DELETE` | `/api/v1/users/{id}` | Delete user (soft) |

### Create user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","role":"user"}'
```

### List users (with pagination + search)

```bash
curl "http://localhost:8080/api/v1/users?page=1&per_page=20&search=alice"
```

### Update user (partial)

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated"}'
```

### Delete user

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## Response Envelopes

### Success (single)
```json
{ "id": 1, "name": "Alice", "email": "alice@example.com", "role": "user", "created_at": "...", "updated_at": "..." }
```

### Success (list)
```json
{
  "data": [...],
  "total": 42,
  "page": 1,
  "per_page": 20,
  "total_pages": 3
}
```

### Error
```json
{ "code": 422, "message": "validation failed", "details": { "email": "must be a valid email address" } }
```

---

## Makefile Commands

```bash
make build           # compile binary → bin/go-user-crud
make run             # run server (source .env first)
make watch           # live reload with air
make test            # unit tests with -race
make test/cover      # coverage report
make test/integration # integration tests (requires TEST_DATABASE_URL)
make lint            # golangci-lint
make fmt             # gofmt + goimports
make check           # fmt + vet + lint + test
make migrate/up      # apply migrations
make migrate/down    # rollback last migration
make migrate/create NAME=add_posts  # new migration
make docker/up       # start docker compose
make docker/down     # stop docker compose
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | **required** | Postgres DSN |
| `PORT` | `8080` | HTTP port |
| `LOG_LEVEL` | `info` | `debug\|info\|warn\|error` |
| `LOG_FORMAT` | `json` | `json\|text` |
| `ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated CORS origins |
| `RATE_LIMIT_RPS` | `100` | Requests per second per IP |

---

## Design Decisions

**No router framework** — Go 1.22 added method+path patterns (`GET /users/{id}`) to the stdlib `http.ServeMux`. No need for chi, gorilla/mux, gin, etc.

**No ORM** — Raw SQL via `pgx/v5`. Explicit queries are readable, debuggable, and have no magic. The repository pattern means you could swap in sqlc-generated code trivially.

**Soft deletes** — `deleted_at` column; rows are never hard-deleted. Add a background job or admin endpoint to purge old records if needed.

**`log/slog`** — Standard library structured logging, available since Go 1.21. Zero external logging dependency.

**Layered architecture** — `domain` has zero external deps. `service` depends only on the `domain` interface. `handler` depends on the `service` interface. Each layer is independently testable via mocks.
