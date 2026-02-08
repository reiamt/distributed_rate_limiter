# Distributed Rate Limiter

A production-style distributed rate limiter written in Go, backed by Redis. Implements per-IP rate limiting with two swappable algorithms, atomic Lua scripting, and standard HTTP rate limit headers.

## Features

- **Sliding window algorithm** (Redis) — atomic Lua script using sorted sets, eliminates fixed-window boundary bursts
- **Token bucket algorithm** (in-memory) — supports burst tolerance with configurable capacity and refill rate
- **Swappable backends** — `Limiter` interface allows switching between Redis and in-memory via env var
- **HTTP middleware** — extracts client IP, enforces limits, returns `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset` headers
- **Graceful shutdown** — signal handling with in-flight request draining and Redis connection cleanup
- **Dockerized** — multi-stage build, Docker Compose with Redis, one command to run

## Architecture

```
cmd/server/           Entry point, config loading, signal handling
internal/
├── config/           Environment-based configuration (godotenv)
├── limiter/
│   ├── limiter.go    Limiter interface + Result struct
│   ├── bucket.go     Token bucket algorithm + Manager with cleanup goroutine
│   └── redis_manager.go  Redis sliding window via Lua script
└── middleware/       HTTP rate limiting middleware
```

Core design: algorithm (`limiter/`) is separated from transport (`middleware/`). The `Limiter` interface decouples the two, enabling the backend to be swapped without changing the middleware.

## Algorithms

**Sliding Window (Redis):** Stores each request timestamp in a sorted set. On every request, prunes entries outside the window, adds the new entry, and counts — all atomically via a Lua script. Prevents the boundary-burst problem inherent in fixed-window counters.

**Token Bucket (in-memory):** Bucket fills at rate `r` tokens/sec up to capacity `b`. Each request consumes one token. Allows bursts up to `b` after idle periods, then throttles to `r` req/sec. Uses double-checked locking for concurrent bucket creation. A background goroutine cleans up idle buckets to prevent memory leaks.

## Quick Start

```bash
docker compose up --build
```

Or run locally:

```bash
# set env vars in .env or export them
go run ./cmd/server
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `MODE` | `redis` | `redis` or `local` |
| `REDIS_ADDR` | `localhost:6379` | Redis connection address |
| `RATELIMIT` | `5` | Max requests per window |
| `PORT` | `:8080` | Server listen port |

## Tests

```bash
go test ./...
```

Unit tests cover token bucket logic (capacity, refill, ceiling), manager (per-key isolation, pointer identity), and middleware (status codes, IP extraction) using table-driven tests and interface mocking.

## Concurrency Patterns Used

- `sync.Mutex` for token bucket state
- `sync.RWMutex` with double-checked locking for bucket map access
- Goroutine + `time.Ticker` + `select` for background cleanup
- `signal.NotifyContext` for graceful shutdown
- Atomic Redis operations via Lua scripting
