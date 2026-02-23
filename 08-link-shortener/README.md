# Task 08 — Link Shortener

## What You'll Learn

- How to combine everything from the previous exercises into a real application
- What a **map** is and how to use it for lookups
- What a **mutex** is and why concurrent access needs protection
- How to define **methods on a struct**
- How to use **path parameters** (`/r/{key}`)
- How to do **HTTP redirects**
- How to use `os.ReadFile` and `os.WriteFile` for simple whole-file I/O

## Key Concepts

### What is a Map?

A **map** is like a phone book. You look up a name (the **key**) and get back
a phone number (the **value**). In Go:

```go
phonebook := map[string]string{
    "Alice": "07700 900001",
    "Bob":   "07700 900002",
}

fmt.Println(phonebook["Alice"]) // prints: 07700 900001
```

The type `map[string]string` means "a map where both keys and values are
strings". You can also check if a key exists:

```go
number, ok := phonebook["Charlie"]
if !ok {
    fmt.Println("Charlie not found")
}
```

The second return value `ok` is `true` if the key exists and `false` if it
doesn't. This is how Go handles "not found" without crashing.

In this exercise, your map will store short keys and their corresponding URLs:
`"0" → "https://go.dev"`, `"1" → "https://github.com"`, etc.

### What is a Mutex?

Imagine a single bathroom in an office. If two people try to use it at the same
time, chaos ensues. So there's a lock on the door — one person locks it, does
their thing, then unlocks it. Everyone else waits their turn.

A **mutex** (short for "mutual exclusion") is the same concept for your code.
When multiple goroutines try to read and write the same data at the same time,
things go wrong. A mutex ensures only one goroutine can access the data at a
time:

```go
var mu sync.Mutex

mu.Lock()
// only one goroutine can be here at a time
data["key"] = "value"
mu.Unlock()
```

In this exercise, your link store uses a mutex to protect the map from
concurrent access. Every time you read or write the map, you lock first and
unlock after.

### Methods on a Struct

So far, your functions have been standalone: `FetchMessage(url)`,
`ReadLines(path)`, etc. But sometimes it makes sense for a function to
**belong** to a type. These are called **methods**.

```go
type LinkStore struct {
    links map[string]string
}

func (s *LinkStore) Shorten(url string) string {
    // s is the LinkStore — like "self" or "this" in other languages
    key := "abc"
    s.links[key] = url
    return key
}
```

The `(s *LinkStore)` part before the function name is the **receiver**. It
means "this function belongs to LinkStore, and `s` is the specific instance".
You call it like `store.Shorten("https://go.dev")`.

### What is an HTTP Redirect?

When you visit a shortened link like `bit.ly/xyz`, the server doesn't send you
a web page. Instead, it sends a **redirect** — a response that says "the thing
you want is actually over here, go there instead". Your browser automatically
follows it.

It's like mail forwarding: you send a letter to an old address, and the post
office sends it on to the new one.

In Go:

```go
http.Redirect(w, r, "https://go.dev", http.StatusFound)
```

`http.StatusFound` is status code `302`, which means "found, but go to this
other URL".

### Path Parameters

In Task 02, your routes were fixed: `/message` always meant the same thing.
But for a link shortener, you need **dynamic** routes: `/r/0`, `/r/1`,
`/r/abc`, etc. The part after `/r/` changes.

Go 1.22 introduced path parameters using curly braces:

```go
mux.HandleFunc("GET /r/{key}", handler)
```

Inside your handler, you extract the value with:

```go
key := r.PathValue("key")
```

If someone visits `/r/42`, then `key` will be `"42"`.

### `os.ReadFile` and `os.WriteFile`

In Task 06 you used `bufio.Scanner` to read line by line — useful for large
files or when you need to process lines individually. But sometimes you just
want to read or write an **entire file** in one go:

```go
data, err := os.ReadFile("links.json")     // read whole file into []byte
err = os.WriteFile("links.json", data, 0644)  // write []byte to file
```

These are simpler than `os.Open` + scanner for cases where the whole file fits
in memory (which is almost always true for config and data files).

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Define a `LinkStore` struct with:
   - `mu sync.Mutex` — protects the map from concurrent access
   - `links map[string]string` — maps short keys to URLs
   - `next int` — counter for generating keys
   - `file string` — path to the JSON file for persistence

3. Write a constructor `NewLinkStore(file string) *LinkStore` that:
   - Returns a new `LinkStore` with an initialised map

4. Write these methods on `*LinkStore`:

   - `Shorten(url string) string` — lock, generate a key using `fmt.Sprintf`
     with `s.next`, store the URL, increment `next`, unlock, return the key

   - `Lookup(key string) (string, bool)` — lock, look up the key in the map,
     unlock, return the URL and whether it was found

   - `SaveToFile() error` — lock, marshal the map to JSON with
     `json.MarshalIndent`, write to `s.file` with `os.WriteFile`, unlock

   - `LoadFromFile() error` — lock, read `s.file` with `os.ReadFile`, unmarshal
     into the map, set `next = len(s.links)`, unlock. If the file doesn't exist
     (`os.IsNotExist`), that's fine — just return nil

5. Define request/response structs:
   - `ShortenRequest` with a `URL` field (JSON tag `"url"`)
   - `ShortenResponse` with `Key` and `ShortURL` fields (JSON tags `"key"`,
     `"short_url"`)

6. Write `ShortenHandler` as a method on `*LinkStore`:
   - Decode the request body
   - If the URL is empty, return 400
   - Call `Shorten`, then `SaveToFile`
   - Respond with status `201` and the key + short URL

7. Write `RedirectHandler` as a method on `*LinkStore`:
   - Extract the key from the path with `r.PathValue("key")`
   - Look up the key — if not found, return 404
   - Redirect to the URL with `http.Redirect` (status 302)

8. In `main`:
   - Create a new `LinkStore` with `"links.json"`
   - Call `LoadFromFile`
   - Set up routes with `http.NewServeMux()`:
     - `"POST /shorten"` → `ShortenHandler`
     - `"GET /r/{key}"` → `RedirectHandler`
   - Start the server on port `8080`

## Run It

```bash
go run .
```

Then in another terminal:

```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://go.dev"}'

# Follow the short link (curl follows redirects by default with -L)
curl -L http://localhost:8080/r/0
```

## Test It

Copy `main_test.go` from the `answers/` directory into this directory, then:

```bash
go test -v
```

The tests check five things:
1. **ShortenHandler** — POST a URL, verify 201 and key in response
2. **RedirectHandler** — pre-populate store, GET `/r/{key}`, verify 302 redirect
3. **Redirect not found** — GET `/r/nonexistent`, verify 404
4. **Shorten empty URL** — POST `{"url":""}`, verify 400
5. **Persistence** — save to file, create new store, load, verify data survives

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Initialise a map: `make(map[string]string)`
- Generate a key: `fmt.Sprintf("%d", s.next)`
- Lock/unlock: `s.mu.Lock()` and `s.mu.Unlock()` (or `defer s.mu.Unlock()`)
- Marshal JSON: `json.MarshalIndent(s.links, "", "  ")`
- Write a file: `os.WriteFile(path, data, 0644)`
- Read a file: `os.ReadFile(path)`
- Check "file not found": `os.IsNotExist(err)`
- Path parameter: `r.PathValue("key")`
- Redirect: `http.Redirect(w, r, url, http.StatusFound)`
- Set status 201: `w.WriteHeader(http.StatusCreated)`
- The imports you'll need: `"encoding/json"`, `"fmt"`, `"net/http"`, `"os"`,
  `"sync"`
