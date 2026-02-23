# CLAUDE.md — Project State & Student Progress

## Project Overview

Go learning repository for a work experience student. Ten progressive exercises
using **only the Go standard library** (plus Docker and Docker Compose for three
exercises, and `go-redis` for the capstone).

## Module

`github.com/student-dev/go-starter` — Go 1.26, macOS ARM64.

## Exercise Map

| # | Directory              | Concepts                              | Status  |
|---|------------------------|---------------------------------------|---------|
| 1 | `01-hello-world`       | `package main`, `import`, `fmt`       | Ready   |
| 2 | `02-local-mock-server` | `net/http`, `net/http/httptest`, JSON  | Ready   |
| 3 | `03-http-client`       | HTTP client, JSON decoding, `io`      | Ready   |
| 4 | `04-concurrency`       | Goroutines, channels, `sync.WaitGroup` | Ready   |
| 5 | `05-docker`            | Dockerfile, `docker build/run`, ports  | Ready   |
| 6 | `06-file-io`           | `os.Open`, `bufio.Scanner`, `os.Create`, `strings.Contains` | Ready   |
| 7 | `07-rest-post`         | POST, `json.NewDecoder/Encoder`, `http.NewServeMux` patterns | Ready   |
| 8 | `08-link-shortener`    | Maps, `sync.Mutex`, persistence, redirects, path params | Ready   |
| 9 | `09-docker-compose`    | Docker Compose, replicas, Nginx, `os.Hostname()` | Ready   |
| 10 | `10-compose-link-shortener` | Redis, `go-redis`, Docker Compose, replicas, Nginx, `os.Hostname()` | Ready   |

## How to Run Tests

```bash
go test ./...
```

## Student Progress

- [ ] Task 01 — Hello World completed
- [ ] Task 02 — Mock Server completed
- [ ] Task 03 — HTTP Client completed
- [ ] Task 04 — Concurrency completed
- [ ] Task 05 — Docker completed
- [ ] Task 06 — File I/O completed
- [ ] Task 07 — REST API with POST completed
- [ ] Task 08 — Capstone Link Shortener completed
- [ ] Task 09 — Docker Compose completed
- [ ] Task 10 — Capstone Link Shortener with Redis completed

## Memory

Key decisions and context for this project:

- Standard library only — no third-party packages (exception: Task 10 uses `go-redis`).
- Tests are the source of truth; exercises 01–08 each have a `main_test.go`.
  Exercises 05, 09, 10 use shell scripts for testing (Docker exercises).
- The mock server from Task 02 is reused in Task 03's tests to validate the client.
- Task 06 teaches file I/O — the fundamental gap before building the capstone.
- Task 07 introduces POST, stream-based JSON, and ServeMux method patterns.
- Task 08 is the capstone — combines maps, mutexes, file persistence, handlers, and redirects.
- Task 09 introduces Docker Compose with replicas and Nginx load balancing.
- Task 10 is the final capstone — combines Redis, Docker Compose, replicas, Nginx, HTTP handlers, and redirects.

## Lessons Learned

Record any logic errors, unexpected behaviour, or corrections here to prevent
future regressions.

- All 4 tests (1 in Task 01, 1 in Task 02, 2 in Task 03) pass on initial setup.
- Task 03's test uses a local `mockMessageHandler` that replicates Task 02's
  handler, keeping the test self-contained (no cross-package imports needed).
- `TestFetchMessageBadURL` verifies the client handles connection errors gracefully.
- Task 04 has 3 tests (correctness, concurrency timing, empty input) — total now 7.
- Task 04's timing test asserts 5×200ms tasks complete in <500ms to prove concurrency.
- Task 05 uses a shell script for testing instead of `main_test.go` (Docker exercises
  don't fit Go's test framework). The `05-docker/answers/` directory has its own
  `go.mod` so the Docker build context is self-contained.
- Task 06 has 4 tests (read, filter, write, not-found) — total now 11.
- Task 07 has 4 tests (GET message, POST greet, empty name, bad JSON) — total now 15.
- Task 07 uses `http.NewServeMux()` with Go 1.22+ method patterns (`"GET /message"`, `"POST /greet"`).
- Task 08 has 5 tests (shorten, redirect, not-found, empty URL, persistence) — total now 20.
- Task 08's redirect test uses `http.Client{CheckRedirect: ...}` to prevent following the redirect.
- Task 08's `LoadFromFile` sets `next = len(links)` so new keys don't collide with loaded ones.
- 20 Go tests (1+1+2+3+4+4+5) across exercises 01–08, plus 3 shell-script-tested exercises (05, 09, 10).
- Task 09 uses Docker Compose with 3 replicas + Nginx; test script verifies multiple hostnames appear.
- Task 09's server uses `os.Hostname()` — inside containers, this returns the container ID.
- Task 09 uses `expose` (internal only) for web replicas and `ports` (host-mapped) for Nginx.
- Docker DNS resolves the service name `web` to all replica IPs for Nginx upstream round-robin.
- Task 10 uses `go-redis` (`github.com/redis/go-redis/v9`) — the first third-party dependency.
- Task 10's `INCR "link:next"` starts at 1 (not 0 like the old counter) — keys are "1", "2", "3".
- Task 10 needs no mutex — Redis handles concurrency atomically.
- Task 10's persistence test restarts only the app container (not Redis) to verify data survives.
- Task 10's Dockerfile copies `go.mod`+`go.sum` before code for dependency layer caching.
- `json.NewEncoder(w).Encode()` adds a trailing `\n` — tests use `json.NewDecoder` which handles this.
- Task 10 now uses 3 replicas + Nginx load balancing (same pattern as Task 09).
- Task 10's `GET /health` endpoint returns `{"status":"ok","hostname":"<id>"}` for proving round-robin.
- Task 10's `ShortenResponse` includes a `hostname` field showing which replica handled the request.
- Task 10 has 6 shell-script tests (shorten, redirect, not-found, empty URL, persistence, load balancing).
