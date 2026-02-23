# Go Starter

A hands-on introduction to Go for beginners. Ten progressive exercises that
build on each other — all using Go's standard library (plus Docker and Docker
Compose for three exercises).

## Prerequisites

- macOS with [Homebrew](https://brew.sh) installed
- A terminal (Terminal.app or iTerm2)
- [Visual Studio Code](https://code.visualstudio.com)

## 1. Install Go

```bash
brew install go
```

Verify the installation:

```bash
go version
```

You should see something like `go version go1.26.0 darwin/arm64`.

## 2. Set Up VSCode

1. Open VSCode.
2. Go to **Extensions** (Cmd+Shift+X).
3. Search for **Go** and install the official extension by the Go team.
4. When prompted, click **Install All** to get `gopls` and other Go tools.
5. Enable **format on save**:
   - Open Settings (Cmd+,)
   - Search for `format on save`
   - Tick the checkbox for **Editor: Format On Save**

This means every time you save a `.go` file, it will automatically be formatted
to standard Go style. No arguments about tabs vs spaces — Go decides for you.

## 3. Clone and Run

```bash
git clone https://github.com/student-dev/go-starter.git
cd go-starter
```

Run all tests to check everything works:

```bash
go test ./...
```

## 4. Exercises

Work through these in order. Each has its own `README.md` with instructions.

| # | Directory | What You'll Learn |
|---|-----------|-------------------|
| 1 | [01-hello-world](./01-hello-world/) | The basics: packages, imports, printing |
| 2 | [02-local-mock-server](./02-local-mock-server/) | Building an HTTP server that returns JSON |
| 3 | [03-http-client](./03-http-client/) | Calling an API and parsing the response |
| 4 | [04-concurrency](./04-concurrency/) | Goroutines, channels, and WaitGroups |
| 5 | [05-docker](./05-docker/) | Containerising your server with Docker |
| 6 | [06-file-io](./06-file-io/) | Reading and writing files, filtering text |
| 7 | [07-rest-post](./07-rest-post/) | POST requests, request bodies, JSON encoding |
| 8 | [08-link-shortener](./08-link-shortener/) | Capstone: maps, mutexes, persistence, redirects |
| 9 | [09-docker-compose](./09-docker-compose/) | Docker Compose, replicas, Nginx load balancing |
| 10 | [10-compose-link-shortener](./10-compose-link-shortener/) | Capstone: Redis, replicas, Nginx, Docker Compose |

## Tips

- **Read the test first.** Exercises 01–08 each have a `main_test.go` that shows
  exactly what your code needs to do. Exercises 05, 09, and 10 use shell scripts
  for testing instead (Docker exercises don't fit Go's test framework).
- **Use `go test -v`** in an exercise directory to see detailed output.
- **Test POST endpoints with curl.** For exercises with POST routes, use
  `curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' URL`.
- **Ask questions.** The goal is to learn, not to finish fast.
