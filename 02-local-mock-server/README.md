# Task 02 — Local Mock Server

## What You'll Learn

- How to build an HTTP server with `net/http`
- How to return JSON from an endpoint
- How to use `net/http/httptest` to test a server without starting it for real
- What **JSON**, **HTTP**, and **servers** actually are

## Key Concepts

### What is a Server?

A **server** is a program that waits for requests and sends back responses.
When you visit a website, your browser (the **client**) sends a request to a
server, and the server sends back the web page.

In this exercise, you'll build a tiny server that runs on your own computer.
Instead of serving web pages, it will send back data.

### What is HTTP?

**HTTP** (HyperText Transfer Protocol) is the system of rules that browsers and
servers use to talk to each other. When you type a URL into your browser, it
sends an **HTTP request**. The server processes it and sends back an **HTTP
response**.

An HTTP response includes:
- A **status code** (like `200` for success, `404` for "not found")
- **Headers** — metadata about the response (like what type of data it contains)
- A **body** — the actual content (the web page, the data, etc.)

### What is JSON?

**JSON** (JavaScript Object Notation) is a way of writing data as text so that
programs can easily read and share it. It looks like this:

```json
{
  "message": "Hello from the backend!"
}
```

Think of it as a list of **key-value pairs**:
- The **key** is a label (like `"message"`) — always in double quotes
- The **value** is the data (like `"Hello from the backend!"`)
- Keys and values are separated by a colon `:`
- The whole thing is wrapped in curly braces `{ }`

JSON can hold different types of data:

```json
{
  "name": "Alice",
  "age": 17,
  "isStudent": true,
  "hobbies": ["coding", "music", "football"]
}
```

Almost every web API in the world uses JSON to send data back and forth. When
your phone app loads your Instagram feed or checks the weather, it's receiving
JSON from a server behind the scenes.

### What is a Handler?

In Go, an HTTP **handler** is a function that processes an incoming request and
writes a response. It always has this signature:

```go
func SomeName(w http.ResponseWriter, r *http.Request)
```

- `w` (ResponseWriter) — you write your response into this. Think of it as a
  blank page you fill in and send back.
- `r` (Request) — this contains everything about the incoming request (the URL,
  headers, body, etc.). You read from it to understand what the client wants.

### What is a Header?

HTTP headers are pieces of extra information attached to a request or response.
The `Content-Type` header tells the client what kind of data the response
contains. Setting it to `application/json` means "this response body is JSON".

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Write a function called `MessageHandler` that:
   - Takes an `http.ResponseWriter` and an `*http.Request`
   - Sets the `Content-Type` header to `application/json`
   - Writes this JSON body: `{"message": "Hello from the backend!"}`

3. In your `main` function, register `MessageHandler` at the path `/message`
   and start a server on port `8080`.

## Run It

```bash
go run .
```

Then open a **new terminal tab** and test it with `curl` (a command-line tool
for making HTTP requests):

```bash
curl http://localhost:8080/message
```

You should see: `{"message": "Hello from the backend!"}`

`localhost` means "this computer" and `:8080` is the port number — like a door
number on a building. Your server is listening on door 8080.

## Test It

Copy `main_test.go` from the `answers/` directory into this directory, then:

```bash
go test -v
```

The test uses `httptest.NewServer` to spin up your handler without needing a
real port. This is how Go developers test HTTP servers — you don't need the
server actually running.

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Set a header: `w.Header().Set("Content-Type", "application/json")`
- Write a response body: `fmt.Fprint(w, "some string")`
- Register a handler: `http.HandleFunc("/message", MessageHandler)`
- Start the server: `http.ListenAndServe(":8080", nil)`
- Your handler function signature: `func MessageHandler(w http.ResponseWriter, r *http.Request)`
