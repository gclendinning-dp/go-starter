# Coming from Python

If you've written Python before, most of Go will feel familiar in shape — it's
the compiler that will surprise you. Here's a quick comparison to help you map
what you already know.

## Compiled and statically typed

This is the biggest change. Python decides what a variable is while your program
runs; Go decides before it runs, and refuses to build if anything doesn't line
up.

```python
# Python
x = 5
x = "now a string"   # perfectly fine
```

```go
// Go
x := 5
x = "now a string" // compile error: cannot use "now a string" as int
```

`x := 5` is Go's **short declaration** — it creates the variable and works out
its type from the value. After that, the type never changes.

Two things trip up nearly everyone coming from Python. First, unused variables
and imports are **errors**, not warnings:

```go
func main() {
    name := "Alice" // compile error: declared and not used
}
```

Second, there's no REPL-driven workflow and no top-to-bottom script execution.
A Go program has one entry point — `func main()` — and you build and run the
whole thing with `go run .`.

The trade-off is worth it: a large class of bug you'd normally find by running
the code (typos in attribute names, wrong argument types, forgetting a return)
gets caught before the program starts.

## No classes

In Python you organise code into classes with `__init__` and methods. Go doesn't
have classes at all. Instead you define **structs** (plain data) and attach
functions to them:

```python
# Python
class Greeter:
    def __init__(self, name):
        self.name = name

    def greet(self):
        return f"Hello, {self.name}!"
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
to a type. It plays the same role as `self` in Python, but it goes before the
function name and you choose what to call it.

There's no `__init__` either. You either fill in the fields directly
(`Greeter{Name: "Alice"}`) or write a plain function that returns one — by
convention, `NewGreeter`.

## No inheritance

In Python you might write `class Dog(Animal)` to inherit behaviour. Go doesn't
have inheritance. Instead, you use **composition** — embedding one struct inside
another:

```python
# Python
class Animal:
    def __init__(self, name):
        self.name = name

class Dog(Animal):
    def bark(self):
        return "Woof!"
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
"a kind of" Animal — it's a Dog that contains an Animal. No MRO, no `super()`,
no diamond problem.

## Interfaces are duck typing, checked by the compiler

You already know duck typing: if it has a `.greet()` method, you can call
`.greet()` on it, and you find out whether you were right when the line runs.

Go works the same way, except it checks at compile time. You never declare that
a type implements an interface — if it has the right methods, it just does. This
is called **structural typing**:

```python
# Python — checked when the line runs
class Greeter:
    def greet(self):
        return "Hi!"

def welcome(g):
    print(g.greet())   # AttributeError at runtime if g has no greet()
```

```go
// Go — checked when you compile
type Greeter interface { Greet() string }

type MyGreeter struct{}
func (g MyGreeter) Greet() string { return "Hi!" }

func Welcome(g Greeter) { fmt.Println(g.Greet()) }
// Passing anything without a Greet() method won't build.
```

So the mental model you already have is right. Go just moves the moment of truth
from "when this line executes in production" to "when I hit save".

## Errors, not exceptions

Python raises exceptions and you catch them. Go doesn't have exceptions.
Instead, functions return an error value alongside their result, and you check
it explicitly:

```python
# Python
try:
    with open("file.txt") as f:
        data = f.read()
except OSError as e:
    print(e)
```

```go
// Go
data, err := os.ReadFile("file.txt")
if err != nil {
    fmt.Println(err)
    return
}
```

Go functions can return more than one value — that's what `data, err :=` is
doing. The `if err != nil` pattern then appears everywhere. It feels repetitive
at first, but it means errors are handled where they happen, and you can see
every failure path by reading straight down the page. No `try` blocks wrapping
half a function, and nothing propagating silently three frames up.

## Concurrency without async/await

Python gives you `threading` (limited by the GIL for CPU work), `asyncio` with
`async`/`await`, and `multiprocessing`. Go has one answer: **goroutines**.

```python
# Python — asyncio
async def fetch(url):
    ...

results = await asyncio.gather(fetch(a), fetch(b))
```

```go
// Go — put "go" in front of a call
go fetch(a)
go fetch(b)
```

There are no `async` or `await` keywords, and no separate "async" versions of
library functions to hunt for. Any function can run in a goroutine. Because
there's no GIL, goroutines doing CPU work genuinely run on multiple cores.

You'll build this properly in [Task 04](../04-concurrency/), including
**channels** — the way goroutines hand results back to each other.

## Other differences you'll notice

| Python | Go |
|--------|-----|
| Indentation defines blocks | Curly braces `{ }` define blocks |
| `import x` / `from x import y` | `package` / `import` |
| `_private` by convention | Uppercase = exported, lowercase = unexported (enforced) |
| `x = 5` | `x := 5` (short declaration) |
| `pip`, `requirements.txt`, `venv` | `go get` + `go.mod` |
| `print()` | `fmt.Println()` |
| f-string `f"Hi {name}"` | `fmt.Sprintf("Hi %s", name)` |
| `list` | slice — `[]string` |
| `dict` | map — `map[string]string` |
| `None` | `nil` (and zero values — see below) |
| `if items:` (truthiness) | `if len(items) > 0` (always explicit) |
| List comprehensions | Plain `for` loops |
| `asyncio` / threads | Goroutines + channels (Task 04) |
| `black` / `PEP 8` arguments | `gofmt`, built in and non-negotiable |

One more thing worth knowing early: Go has no `None` for ordinary values. Every
type has a **zero value** it starts at — `""` for strings, `0` for numbers,
`false` for booleans. So `var name string` is an empty string, not `None`, and
you don't need the `if x is None` guards you're used to writing.

Don't worry about memorising all of this — it'll make more sense as you work
through the exercises. The key thing: Go is deliberately simple. There's usually
one way to do something, not five.

---

Ready to start? Head to [Task 01 — Hello World](../01-hello-world/).
