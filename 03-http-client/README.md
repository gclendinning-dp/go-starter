# Task 03 — HTTP Client

## What You'll Learn

- How to make HTTP requests with `net/http`
- How to read a response body with `io`
- How to parse JSON into a Go **struct** with `encoding/json`
- What **structs**, **JSON parsing**, and **error handling** are

## Key Concepts

### Client vs Server

In Task 02 you built a **server** — a program that waits for requests and sends
responses. Now you're building a **client** — a program that *sends* requests
and reads the responses.

This is how most software works: clients and servers talking to each other. Your
web browser is a client. Mobile apps are clients. And now you're writing one
from scratch in Go.

### What is a Struct?

A **struct** is Go's way of grouping related data together. Think of it like a
form with labelled fields:

```go
type Person struct {
    Name string
    Age  int
}
```

This defines a new type called `Person` with two fields. You can create one and
access its fields with a dot:

```go
p := Person{Name: "Alice", Age: 17}
fmt.Println(p.Name) // prints: Alice
fmt.Println(p.Age)  // prints: 17
```

Structs are used everywhere in Go. They're how you model real things in your
code — a user, a message, an HTTP response, a database row.

### What is JSON Parsing (Unmarshalling)?

In Task 02, your server sent JSON as plain text:

```json
{"message": "Hello from the backend!"}
```

But a Go program can't do much with raw text. You need to convert it into a Go
value you can work with — like a struct. This conversion is called **parsing**
or **unmarshalling**.

You define a struct that matches the shape of the JSON:

```go
type Response struct {
    Message string `json:"message"`
}
```

The `` `json:"message"` `` part is called a **struct tag**. It tells Go's JSON
parser: "when you see a JSON key called `message`, put its value into this
`Message` field." The tag maps between the JSON world (lowercase `message`) and
the Go world (uppercase `Message`).

Then you parse:

```go
var result Response
json.Unmarshal(body, &result)
// Now result.Message contains "Hello from the backend!"
```

The `&` means "here's the address of `result`" — it lets `Unmarshal` fill in
the struct's fields directly. (This is called a **pointer** — you'll learn more
about these later.)

### What is Error Handling?

Things can go wrong: the server might be offline, the network might drop, the
response might not be valid JSON. Go handles this by returning an **error** value
alongside the result:

```go
resp, err := http.Get(url)
```

This returns two things: the response (`resp`) and an error (`err`). If
everything went fine, `err` is `nil` (Go's word for "nothing"). If something
went wrong, `err` contains a description of the problem.

The pattern in Go is always the same — check the error immediately:

```go
if err != nil {
    // something went wrong, handle it
    return "", err
}
// everything's fine, carry on
```

You'll see this pattern hundreds of times in Go code. It's verbose, but it
forces you to think about what can go wrong — which makes programs more
reliable.

### What is `defer`?

When you open an HTTP response, you need to close it when you're done (like
closing a file). `defer` tells Go "run this line later, when the function ends":

```go
resp, err := http.Get(url)
if err != nil {
    return "", err
}
defer resp.Body.Close() // this runs when the function returns
```

It's a safety net — even if something goes wrong later in the function, the
body will always get closed.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Define a struct called `Response` with one field:
   - `Message string` — tagged with `` `json:"message"` `` so Go knows how to
     map the JSON key to the struct field.

3. Write a function called `FetchMessage` that:
   - Takes a `url` string as its argument
   - Makes an HTTP GET request to that URL
   - Reads the response body
   - Parses the JSON into a `Response` struct
   - Returns the `Message` string and an `error`

4. In your `main` function, call `FetchMessage` with
   `"http://localhost:8080/message"` and print the result.

   To test this manually, start the server from Task 02 first:
   ```bash
   cd ../02-local-mock-server/answers && go run . &
   cd ../../03-http-client && go run .
   ```

## Test It

Copy `main_test.go` from the `answers/` directory into this directory, then:

```bash
go test -v
```

The test spins up a mock version of the Task 02 server automatically — you
don't need to start it yourself.

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Define the struct:
  ```go
  type Response struct {
      Message string `json:"message"`
  }
  ```
- Make a GET request: `resp, err := http.Get(url)`
- Read the body: `body, err := io.ReadAll(resp.Body)`
- Parse JSON: `json.Unmarshal(body, &result)`
- Always close the body: `defer resp.Body.Close()`
- Return multiple values: `func FetchMessage(url string) (string, error)`
