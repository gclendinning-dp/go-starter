# Coming from C#

If you've written C# before, Go will feel quite different. Here's a quick
comparison to help you map what you already know.

## No classes

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

## No inheritance

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

## Interfaces are implicit

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

## Errors, not exceptions

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

## Other differences you'll notice

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

---

Ready to start? Head to [Task 01 — Hello World](../01-hello-world/).
