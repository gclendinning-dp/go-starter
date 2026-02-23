# CLAUDE.md — Project State & Student Progress

## Project Overview

Go learning repository for a work experience student. Eight progressive exercises
using **only the Go standard library** (plus Docker for one exercise).

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

## Memory

Key decisions and context for this project:

- Standard library only — no third-party packages.
- Tests are the source of truth; each exercise has a `main_test.go`.
- The mock server from Task 02 is reused in Task 03's tests to validate the client.
- Task 06 teaches file I/O — the fundamental gap before building the capstone.
- Task 07 introduces POST, stream-based JSON, and ServeMux method patterns.
- Task 08 is the capstone — combines maps, mutexes, file persistence, handlers, and redirects.

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
- `json.NewEncoder(w).Encode()` adds a trailing `\n` — tests use `json.NewDecoder` which handles this.
