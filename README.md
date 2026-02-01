# Tu Crédito – Credit Decision & Management API

Go backend API for managing clients, banks, and credits. It exposes REST endpoints, uses PostgreSQL and Redis, processes credit creation via a worker pool, publishes domain events (mock publisher, Kafka-ready interface), and includes a rule-based decision engine for eligibility and bank routing.

## Requirements

- **Go 1.23+**
- **PostgreSQL 16** (or compatible)
- **Redis 7** (optional; used for caching and rate limiting)
- **Docker & Docker Compose** (for running the stack)
- **golang-migrate** (for running migrations locally; optional when using Docker)

## Project layout

```
.
├── cmd/server/           # Application entrypoint
├── internal/
│   ├── cache/            # Redis cache (and rate-limit primitives)
│   ├── decision/         # Credit routing & eligibility engine (rules, waterfall)
│   ├── domain/           # Entities and domain events
│   ├── event/            # Event publisher (mock Kafka)
│   ├── handler/          # HTTP handlers (REST)
│   ├── middleware/       # Logging, recovery, rate limit
│   ├── metrics/          # Prometheus-style metrics
│   ├── repository/       # Interfaces and mocks
│   │   └── postgres/     # PostgreSQL implementations
│   ├── server/           # Wiring and HTTP server
│   └── service/          # Business logic (worker pool, events)
├── benchmarks/           # Credit service benchmarks
├── migrations/           # SQL schema (golang-migrate, up/down)
├── pkg/
│   ├── config/           # Env-based config
│   ├── httputil/         # JSON responses
│   └── logger/           # Structured logging (zap)
├── Dockerfile
├── docker-compose.yml
├── docker-compose.redis.local.yml   # Redis only (local dev)
├── makefile
└── .github/workflows/deploy.yml     # CI/CD
```

## Architecture

- **Layers**: Handlers → Services → Repositories; domain and events are separate. Easy to swap persistence or plug in a real Kafka producer.
- **Concurrency**: Credit creation is processed by a **worker pool** (goroutines + channel). Validations and eligibility run inside workers; client and bank lookups can run in parallel.
- **Events**: Domain events (`CreditCreated`, `CreditApproved`, `CreditRejected`) are published via an interface; the current implementation is an in-memory mock. Replacing it with a Kafka producer keeps the same API.
- **Caching**: Credits are cached in Redis by ID (with TTL). Rate limiting uses Redis `INCR` + `EXPIRE` per client IP (100 requests per 60 seconds by default).
- **Decision engine**: Extensible rule-based engine in `internal/decision`. Rules run in order (waterfall); first approval wins. Includes payment-range and bank-type rules; easy to add priority, yield, or inventory logic.
- **Observability**: Structured logging (zap), Prometheus-style metrics at `/metrics`, `/health` (liveness), `/ready` (readiness with Postgres/Redis). pprof at `:6060/debug/pprof/` when `PPROF_ENABLED=true`.

## How to run

### With Docker Compose (recommended)

```bash
docker compose -f docker-compose.yml up -d --build
```

Or using the makefile:

```bash
make compose-up
```

- **API**: `http://localhost:8080`
- **Health**: `http://localhost:8080/health`
- **Ready**: `http://localhost:8080/ready`
- **Metrics**: `http://localhost:8080/metrics`
- **PostgreSQL**: `localhost:5432` (user `postgres`, password `postgres`, DB `tucredito`)
- **Redis**: `localhost:6379`

Migrations run automatically via the `migrate` service before the API starts.

### Local (without Docker)

1. Start PostgreSQL and (optionally) Redis. For Redis only: `make redis-up` (uses `docker-compose.redis.local.yml`).
2. Create the database: `createdb tucredito`
3. Run migrations:

   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"
   make migrations-up
   ```

   Or with `migrate` CLI directly:

   ```bash
   migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable" up
   ```

4. Set environment variables (see `.env.example`; defaults below):

   ```bash
   export HTTP_PORT=8080
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"
   export REDIS_ADDR=localhost:6379
   export REDIS_PASSWORD=
   export REDIS_DB=0
   export LOG_LEVEL=info
   export PPROF_ENABLED=true
   ```

5. Run the server:

   ```bash
   go run ./cmd/server
   ```

   Or: `make go-run`


> Additionally it's possible to debug this repository using a `.vscode` folder and a `launch.json` file, like this one:
```
launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/server/main.go"
        }
    ],
    "compounds": []
}
```

## API overview

Health and metrics (no version prefix):

| Method | Path       | Description                    |
|--------|------------|--------------------------------|
| GET    | `/health`  | Liveness                       |
| GET    | `/ready`   | Readiness (Postgres, Redis)    |
| GET    | `/metrics` | Prometheus-style metrics      |

**Clients** (`/v1/clients`):

| Method | Path                         | Description              |
|--------|------------------------------|--------------------------|
| POST   | `/v1/clients`                | Create client            |
| GET    | `/v1/clients`                | List clients (pagination)|
| GET    | `/v1/clients/{id}`           | Get client               |
| PUT    | `/v1/clients/{id}`           | Update client            |
| DELETE | `/v1/clients/{id}`           | Delete (soft) client     |
| POST   | `/v1/clients/{id}/reenable`  | Re-enable client         |
| GET    | `/v1/clients/{id}/credits`  | List credits for client  |

**Banks** (`/v1/banks`):

| Method | Path                        | Description        |
|--------|-----------------------------|--------------------|
| POST   | `/v1/banks`                 | Create bank        |
| GET    | `/v1/banks`                 | List banks         |
| GET    | `/v1/banks/{id}`            | Get bank           |
| PUT    | `/v1/banks/{id}`            | Update bank        |
| DELETE | `/v1/banks/{id}`            | Delete (soft) bank |
| POST   | `/v1/banks/{id}/reenable`   | Re-enable bank     |

**Credits** (`/v1/credits`):

| Method | Path                        | Description                    |
|--------|-----------------------------|--------------------------------|
| POST   | `/v1/credits`               | Create credit (worker pool, events, cache) |
| GET    | `/v1/credits`               | List credits                   |
| GET    | `/v1/credits/{id}`          | Get credit (cache-first)       |
| PUT    | `/v1/credits/{id}`          | Update credit                  |
| DELETE | `/v1/credits/{id}`          | Delete (soft) credit           |
| POST   | `/v1/credits/{id}/reenable` | Re-enable credit               |

## Postman

There is a entire Postman colletion to test any of these endpoints, you have to import the collection and the environment located in:
- Collection: --------> `./postman/TuCredito.postman_collection.json`
- Environment: -----> `./postman/TuCredito.postman_environment.json` 

## Tests and benchmarks

```bash
# Unit tests (no DB/Redis)
make go-test
# or
go test ./... -v -count=1 -short

# Integration tests (require Postgres; set DATABASE_URL)
go test -tags=integration ./internal/repository/postgres/... -v -count=1

# Benchmarks
make go-bench
# or
go test ./benchmarks/... -bench=. -benchmem -run=^$ -count=1
```

## CI/CD pipeline

The pipeline (`.github/workflows/deploy.yml`) runs on push/PR to `main` or `master`:

1. **Lint**: `golangci-lint`
2. **Test**: `go test ./... -v -count=1 -short`
3. **Benchmarks**: `go test ./benchmarks/... -bench=. -benchmem -run=^$ -count=1`
4. **Build**: `go build -o server ./cmd/server`
5. **Integration**: Postgres service, run migrations (`migrations/*.up.sql`), then `go test -tags=integration ./internal/repository/postgres/...`
6. **Deploy**: Placeholder step (build only)

To run integration tests locally with Docker Postgres:

```bash
docker run -d --name tucredito-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=tucredito -p 5432:5432 postgres:16-alpine
# Wait for ready, then:
for f in migrations/*.up.sql; do PGPASSWORD=postgres psql -h localhost -U postgres -d tucredito -f "$f"; done
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/tucredito?sslmode=disable"
go test -tags=integration ./internal/repository/postgres/... -v -count=1
```

## Makefile targets

| Target            | Description                    |
|-------------------|--------------------------------|
| `compose-up`      | Start full stack (API, Postgres, Redis, migrate) |
| `compose-down`    | Stop and remove volumes        |
| `redis-up` / `redis-down` | Redis only (local dev)  |
| `migrations-up`   | Run migrations (set `DATABASE_URL`) |
| `migrations-down` | Rollback one migration         |
| `create-migration`| Create new migration (set `NAME`)   |
| `go-lint` / `go-lint-fix` | Lint                      |
| `go-test`         | Unit tests                     |
| `go-bench`        | Benchmarks                     |
| `go-build` / `go-run` | Build and run binary       |

## Performance notes

- **Credit creation**: Throughput is bounded by worker pool size (default 10) and DB/Redis latency. Increase pool size or scale replicas for higher load.
- **Rate limiting**: 100 requests per 60 seconds per client (Redis). Ensure Redis has enough memory and connections for your traffic.
- **Metrics**: In-memory counters and duration samples; scrape `/metrics` with Prometheus for production.

## AI Use

This proyect leverage the use of AI in the following scenarious: 
- Complete repetitive code: Like `test-cases`, `structs` and `mocks`.
- Write more explicative comments in some needed sections.
- Auto-complete sections: Tool in some IDEs to auto-complete `variables` and `logic` while coding.
- Parts of this `README.md`: To structure a well description.

> Final note: All the features and code related to AI was thoroughly reviewed and validated (according to the official documentation) to avoid AI hallucinations and errors

