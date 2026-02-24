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

Run the answer tests to check everything works:

```bash
go test ./01-hello-world/answers/ ./02-local-mock-server/answers/ ./03-http-client/answers/ ./04-concurrency/answers/ ./06-file-io/answers/ ./07-rest-post/answers/ ./08-link-shortener/answers/
```

## 4. Coming from C#

If you've written C# before, Go will feel quite different. Here's a quick
comparison to help you map what you already know.

### No classes

C# is object-oriented — you organise code into classes with constructors,
properties, and methods. Go doesn't have classes at all. Instead you define
**structs** (plain data) and attach functions to them:

```csharp
// C#
public class Greeter {
    public string Name { get; set; }
    public string Greet() => $"Hello, {Name}!";
}
```

```go
// Go
type Greeter struct {
    Name string
}

func (g Greeter) Greet() string {
    return fmt.Sprintf("Hello, %s!", g.Name)
}
```

The `(g Greeter)` part is called a **receiver** — it's how Go attaches a method
to a type. It's similar to `this` in C#, but you choose the name yourself.

### No inheritance

In C# you might write `class Dog : Animal` to inherit behaviour. Go doesn't
have inheritance. Instead, you use **composition** — embedding one struct inside
another:

```csharp
// C#
class Animal { public string Name { get; set; } }
class Dog : Animal { public void Bark() { ... } }
```

```go
// Go
type Animal struct { Name string }
type Dog struct {
    Animal          // embedded — Dog "has an" Animal
}
func (d Dog) Bark() string { return "Woof!" }
```

`Dog` gets all of `Animal`'s fields and methods automatically, but it's not
"a kind of" Animal — it's a Dog that contains an Animal.

### Interfaces are implicit

In C# you write `class Foo : IBar` to say "Foo implements IBar". In Go, you
don't declare that a type implements an interface — if it has the right methods,
it just does. This is called **structural typing**:

```csharp
// C#
interface IGreeter { string Greet(); }
class Greeter : IGreeter { public string Greet() => "Hi!"; }
```

```go
// Go
type Greeter interface { Greet() string }

type MyGreeter struct{}
func (g MyGreeter) Greet() string { return "Hi!" }
// MyGreeter implements Greeter automatically — no declaration needed.
```

### Errors, not exceptions

C# uses `try/catch` for error handling. Go doesn't have exceptions. Instead,
functions return an error value that you check explicitly:

```csharp
// C#
try {
    var data = File.ReadAllText("file.txt");
} catch (Exception e) {
    Console.WriteLine(e.Message);
}
```

```go
// Go
data, err := os.ReadFile("file.txt")
if err != nil {
    fmt.Println(err)
    return
}
```

This `if err != nil` pattern appears everywhere in Go. It feels repetitive at
first, but it means errors are always handled where they happen — no surprises
from uncaught exceptions three layers up the call stack.

### Other differences you'll notice

| C# | Go |
|----|-----|
| `namespace` / `using` | `package` / `import` |
| `public` / `private` keywords | Uppercase = exported, lowercase = unexported |
| `var x = 5` or `int x = 5` | `x := 5` (short declaration) |
| NuGet packages | `go get` + `go.mod` |
| `async` / `await` | Goroutines + channels (Task 04) |
| `Console.WriteLine()` | `fmt.Println()` |
| Semicolons required | No semicolons |
| Curly braces on new line (convention) | Opening brace on same line (enforced) |

Don't worry about memorising all of this — it'll make more sense as you work
through the exercises. The key thing: Go is deliberately simple. There's usually
one way to do something, not five.

## 5. Exercises

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
| 8 | [08-link-shortener](./08-link-shortener/) | Maps, mutexes, persistence, redirects |
| 9 | [09-docker-compose](./09-docker-compose/) | Docker Compose, replicas, Nginx load balancing |
| 10 | [10-compose-link-shortener](./10-compose-link-shortener/) | Capstone: Redis, replicas, Nginx, Docker Compose |

## Test-Driven Development (TDD)

This project uses **test-driven development**. Each exercise (01–04, 06–08) has a
`main_test.go` already in the exercise directory — before you write any code.
The tests define exactly what your code needs to do.

The TDD cycle is:

1. **Red** — Read the test. Try running it (`go test -v`). It won't compile yet
   because the functions don't exist. That's expected.
2. **Green** — Write the minimum code in `main.go` to make the test pass.
3. **Refactor** — Clean up your code while keeping the tests green.

This is how professional developers work: the test is the specification. You
don't guess what to build — the test tells you. If the test passes, you're done.

Exercises 05, 09, and 10 use shell scripts for testing instead (Docker exercises
don't fit Go's test framework). These scripts run automatically against your
running containers.

### Running tests

```bash
# Run the test for a specific exercise (from the exercise directory):
cd 01-hello-world
go test -v

# Run all answer tests (to verify the repo setup):
go test ./01-hello-world/answers/ ./02-local-mock-server/answers/ ...
```

## Tips

- **Read the test first.** The `main_test.go` in each exercise directory is your
  specification. Read it before writing any code — it tells you what functions to
  create, what they should return, and what edge cases to handle.
- **Use `go test -v`** for detailed output showing which tests pass and fail.
- **Test POST endpoints with curl.** For exercises with POST routes, use
  `curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' URL`.
- **Ask questions.** The goal is to learn, not to finish fast.
