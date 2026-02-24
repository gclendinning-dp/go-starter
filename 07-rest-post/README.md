# Task 07 — REST API with POST

## What You'll Learn

- The difference between **GET** and **POST** requests
- How to read a **request body** and parse JSON from it
- How to use `json.NewDecoder` and `json.NewEncoder` (stream-based JSON)
- How to use `http.NewServeMux` with **method patterns** for routing
- How to return proper **error responses** with `http.Error`

## Key Concepts

### GET vs POST

In Task 02 you built a server that handled **GET** requests — the client just
asks "give me some data" and the server responds. That's like looking at a
noticeboard: you read what's there.

A **POST** request is different. The client **sends data** to the server, and
the server does something with it. It's like filling in a form and handing it
to someone at a desk — you're providing information, not just reading.

- **GET** = "Give me something" (reading)
- **POST** = "Here's some data, do something with it" (writing/creating)

### What is a Request Body?

When you send a POST request, the data goes in the **body** of the request.
Think of it like an envelope: the URL is the address on the front, and the
body is the letter inside.

In this exercise, the client will send a JSON body like `{"name": "Alice"}`,
and your server will read it and respond with a greeting.

### Stream-Based JSON: `json.NewDecoder` / `json.NewEncoder`

In Task 03, you parsed JSON like this:

```go
body, err := io.ReadAll(resp.Body)       // read everything into memory
err = json.Unmarshal(body, &result)       // then parse it
```

That works, but there's a better way when you're reading from a stream
(like an HTTP request body):

```go
err := json.NewDecoder(r.Body).Decode(&req)   // read and parse in one step
```

And for writing JSON responses:

```go
json.NewEncoder(w).Encode(resp)   // write JSON directly to the response
```

Why is this better? The `ReadAll` + `Unmarshal` approach reads the entire body
into memory first, then parses it. The decoder approach **streams** the data —
it reads and parses at the same time, which is more efficient. For small
payloads the difference doesn't matter, but it's the idiomatic Go way.

### Method Patterns with `http.NewServeMux`

In Task 02 you used `http.HandleFunc("/message", MessageHandler)`, which
matches **any** HTTP method (GET, POST, DELETE, etc.) on that path.

Go 1.22 introduced a better way — you can specify the method in the pattern:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /message", MessageHandler)
mux.HandleFunc("POST /greet", GreetHandler)
```

Now `GET /message` only matches GET requests to `/message`, and
`POST /greet` only matches POST requests to `/greet`. If someone sends a
GET to `/greet`, the server automatically returns `405 Method Not Allowed`.

### Error Responses with `http.Error`

When something goes wrong (bad input, missing data), you should tell the client
what happened. `http.Error` sends a plain-text error message with the right
status code:

```go
http.Error(w, "name is required", http.StatusBadRequest)
return
```

`http.StatusBadRequest` is the constant for status code `400`, which means
"the client sent something the server can't process". Always `return` after
calling `http.Error` so you don't accidentally write more data to the response.

## Before You Start — Read the Test

Open `main_test.go` in this directory. It shows how `MessageHandler` and
`GreetHandler` are tested, including the `GreetResponse` struct the test expects,
and the edge cases for empty names and invalid JSON.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Define two structs:
   - `GreetRequest` with a `Name` field (string, JSON tag `"name"`)
   - `GreetResponse` with a `Greeting` field (string, JSON tag `"greeting"`)

3. Write a function called `MessageHandler` that:
   - Takes an `http.ResponseWriter` and `*http.Request`
   - Sets `Content-Type` to `application/json`
   - Returns `{"message": "Hello from the backend!"}` (same as Task 02)

4. Write a function called `GreetHandler` that:
   - Takes an `http.ResponseWriter` and `*http.Request`
   - Decodes the request body as JSON into a `GreetRequest`
   - If the JSON is invalid, returns a `400` error with `http.Error`
   - If the name is empty, returns a `400` error
   - Otherwise, sets `Content-Type` to `application/json` and responds with
     `{"greeting": "Hello, Alice!"}` (using the name from the request)

5. In your `main` function:
   - Create a `http.NewServeMux()`
   - Register `MessageHandler` at `"GET /message"`
   - Register `GreetHandler` at `"POST /greet"`
   - Start the server on port `8080` using the mux

## Run It

```bash
go run .
```

Then in another terminal, test the GET endpoint:

```bash
curl http://localhost:8080/message
```

And test the POST endpoint:

```bash
curl -X POST http://localhost:8080/greet \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```

You should see: `{"greeting":"Hello, Alice!"}`

## Test It

```bash
go test -v
```

The tests check four things:
1. **MessageHandler** — GET /message returns the correct JSON
2. **GreetHandler** — POST with a valid name returns the greeting
3. **GreetHandler (empty name)** — POST with `{"name":""}` returns 400
4. **GreetHandler (bad JSON)** — POST with invalid JSON returns 400

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Decode request body: `json.NewDecoder(r.Body).Decode(&req)`
- Encode a response: `json.NewEncoder(w).Encode(resp)`
- Return an error: `http.Error(w, "message", http.StatusBadRequest)`
- Create a mux: `mux := http.NewServeMux()`
- Register with a method: `mux.HandleFunc("POST /greet", GreetHandler)`
- Start the server with a mux: `http.ListenAndServe(":8080", mux)`
- Build the greeting string: `fmt.Sprintf("Hello, %s!", req.Name)`
- The imports you'll need: `"encoding/json"`, `"fmt"`, `"net/http"`
