# Task 10 — Capstone: Link Shortener with Redis

## What You'll Learn

- What **Redis** is and why it's useful for storing data
- How to use a **third-party Go package** (`go-redis`)
- How **key-value stores** work
- How to combine Docker Compose, HTTP handlers, and external services
- What `context.Context` is (briefly)
- Redis commands: `INCR`, `HSET`, `HGET`

## Key Concepts

### What is Redis?

Imagine a **whiteboard in a shared office**. Anyone can walk up and write
something on it, read what's there, or erase something. It's fast (no flipping
through filing cabinets), everyone can access it, and the information stays
there until someone removes it.

**Redis** is like that whiteboard, but for your programs. It's a **database**
that stores data in memory (RAM), which makes it extremely fast. Unlike the
file-based persistence from the previous exercises, Redis:

- Handles **concurrent access** for you (no mutexes needed)
- **Survives restarts** (it saves to disk periodically)
- Can be **shared** between multiple services
- Is **much faster** than reading/writing files

### Key-Value Stores

Redis is a **key-value store** — the simplest kind of database. You store data
as pairs: a key (the name) and a value (the data). It's like a Go map:

```
"name"    → "Alice"
"country" → "UK"
"age"     → "25"
```

Redis also supports **hashes**, which are maps within maps:

```
"links" → {
    "1" → "https://go.dev"
    "2" → "https://github.com"
}
```

In this exercise, you'll use a Redis hash called `"links"` to store your
shortened URLs, and a key called `"link:next"` to track the next ID.

### Redis Commands

You'll use three Redis commands:

- **`INCR key`** — increments a number stored at `key` by 1 and returns the new
  value. If the key doesn't exist, it starts at 0 and increments to 1. This is
  **atomic** — even if two requests run at the same time, they'll get different
  numbers.

- **`HSET hash field value`** — sets a field in a hash. Like
  `map["field"] = "value"` in Go.

- **`HGET hash field`** — gets a field from a hash. Like `map["field"]` in Go.

### Third-Party Packages in Go

Every exercise so far has used only Go's standard library. This exercise
introduces your first **third-party package**: `go-redis`.

To add a dependency:

```bash
go get github.com/redis/go-redis/v9
```

This downloads the package and adds it to your `go.mod` and `go.sum` files.
`go.sum` is a lockfile — it records the exact checksums of every dependency to
ensure reproducible builds.

You can also run `go mod tidy` to clean up — it adds missing dependencies and
removes unused ones.

### `context.Context` (Brief)

You'll see `context.Background()` in the Redis client calls. A context carries
deadlines, cancellation signals, and request-scoped values. For now, think of
`context.Background()` as passing a "no special instructions" token — it means
"just do the normal thing, no timeout, no cancellation."

### Why Redis Instead of a File?

In the previous link shortener exercise, you used a file (`links.json`) and a
mutex to handle persistence and concurrency. That works, but has limitations:

| File-based | Redis |
|------------|-------|
| Need a mutex for thread safety | Redis handles concurrency atomically |
| Data lost if you forget to save | Redis persists automatically |
| Only one process can use the file | Multiple services can share Redis |
| Read/write entire file each time | Read/write individual keys |

### Docker Compose for Multiple Services

Building on Task 09, this exercise uses Docker Compose to run **two services**:
Redis and your Go application. The `depends_on` key ensures Redis starts first,
and Docker DNS lets your app connect to Redis using the hostname `redis`.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Add the Redis dependency:

   First create a `go.mod`:

   ```
   module shortener

   go 1.26.0
   ```

   Then add the dependency:

   ```bash
   go get github.com/redis/go-redis/v9
   ```

3. Define a `LinkStore` struct with:
   - `rdb *redis.Client` — the Redis client
   - `ctx context.Context` — the context for Redis operations

4. Write a constructor `NewLinkStore(redisAddr string) *LinkStore` that:
   - Creates a `redis.NewClient` with the given address
   - Sets `ctx` to `context.Background()`

5. Write these methods on `*LinkStore`:

   - `Shorten(url string) (string, error)` — use `INCR "link:next"` to get a
     unique key, then `HSET "links" key url` to store it. Return the key.

   - `Lookup(key string) (string, bool)` — use `HGET "links" key`. If the key
     doesn't exist, return `("", false)`.

6. Define request/response structs (same as before):
   - `ShortenRequest` with a `URL` field (JSON tag `"url"`)
   - `ShortenResponse` with `Key` and `ShortURL` fields

7. Write `ShortenHandler` as a method on `*LinkStore`:
   - Decode the request body
   - If the URL is empty, return 400
   - Call `Shorten` — if it returns an error, return 500
   - Respond with status `201` and the key + short URL

8. Write `RedirectHandler` as a method on `*LinkStore`:
   - Extract the key from the path with `r.PathValue("key")`
   - Look up the key — if not found, return 404
   - Redirect with `http.Redirect` (status 302)

9. In `main`:
   - Read `REDIS_ADDR` from the environment (default `"localhost:6379"`)
   - Create a `NewLinkStore` with that address
   - Set up routes: `"POST /shorten"` and `"GET /r/{key}"`
   - Start the server on port 8080

10. Create a `Dockerfile`:

    ```dockerfile
    FROM golang:1.26-alpine
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY main.go .
    RUN go build -o server .
    EXPOSE 8080
    CMD ["./server"]
    ```

    Notice the two-step copy: `go.mod` and `go.sum` are copied first, then
    `go mod download` runs. This means Docker **caches** the dependency
    download layer — if you change `main.go` but not your dependencies, Docker
    won't re-download them. This is a common optimization.

11. Create a `docker-compose.yml`:

    ```yaml
    services:
      redis:
        image: redis:7-alpine
        ports:
          - "6379:6379"

      app:
        build: .
        ports:
          - "8080:8080"
        environment:
          - REDIS_ADDR=redis:6379
        depends_on:
          - redis
    ```

## Run It

```bash
docker compose up -d --build
```

Then in another terminal:

```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://go.dev"}'

# Follow the short link
curl -L http://localhost:8080/r/1
```

Stop everything:

```bash
docker compose down -v
```

## Test It

An automated test script is provided in the `answers/` directory:

```bash
bash answers/test_shortener.sh
```

The script tests five things:
1. **POST /shorten** — verify 201 and non-empty key
2. **GET /r/{key}** — verify 302 redirect to the correct URL
3. **GET /r/nonexistent** — verify 404
4. **POST /shorten with empty URL** — verify 400
5. **Persistence** — restart the app container, verify data survives in Redis

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Create a Redis client: `redis.NewClient(&redis.Options{Addr: addr})`
- Increment a key: `s.rdb.Incr(s.ctx, "link:next").Result()`
- Set a hash field: `s.rdb.HSet(s.ctx, "links", key, url).Err()`
- Get a hash field: `s.rdb.HGet(s.ctx, "links", key).Result()`
- Read an env var: `os.Getenv("REDIS_ADDR")`
- The imports you'll need: `"context"`, `"encoding/json"`, `"fmt"`, `"net/http"`,
  `"os"`, `"github.com/redis/go-redis/v9"`
