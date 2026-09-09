# Coming from Java

If you've written Java before, Go will feel like someone took the same ideas and
removed about two-thirds of the ceremony. Here's a quick comparison to help you
map what you already know.

## No classes

Java puts everything inside a class — even `main`. Go doesn't have classes at
all. Functions live directly in a **package**, and data lives in **structs**
that you attach functions to:

```java
// Java
public class Greeter {
    private String name;
    public Greeter(String name) { this.name = name; }
    public String greet() { return "Hello, " + name + "!"; }
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
to a type. It's similar to `this` in Java, but you choose the name yourself.

There are no constructors either. You either fill in the fields directly
(`Greeter{Name: "Alice"}`) or write a plain function that returns one — by
convention, `NewGreeter`. And your program's entry point is just
`func main()`, not `public static void main(String[] args)` buried in a class.

## No inheritance

In Java you write `class Dog extends Animal` to inherit behaviour. Go doesn't
have inheritance — no `extends`, no `abstract`, no `super`. Instead, you use
**composition** — embedding one struct inside another:

```java
// Java
class Animal { String name; }
class Dog extends Animal { String bark() { return "Woof!"; } }
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
"a kind of" Animal — it's a Dog that contains an Animal. The class hierarchies
you're used to designing don't have an equivalent here, and that's deliberate.

## Interfaces are implicit

In Java you write `class Foo implements IBar` to say so explicitly. In Go, you
don't declare that a type implements an interface — if it has the right methods,
it just does. This is called **structural typing**:

```java
// Java
interface Greeter { String greet(); }
class MyGreeter implements Greeter {
    public String greet() { return "Hi!"; }
}
```

```go
// Go
type Greeter interface { Greet() string }

type MyGreeter struct{}
func (g MyGreeter) Greet() string { return "Hi!" }
// MyGreeter implements Greeter automatically — no declaration needed.
```

This changes how interfaces get designed. In Java the interface usually comes
first and implementations are written to fit it. In Go it's common to write the
concrete type first, then define a small interface — often one or two methods —
at the place that needs it. You can also define an interface that an existing
type from someone else's package already satisfies, without touching their code.

## Errors, not exceptions

Java has checked exceptions: `throws IOException` forces callers to acknowledge
that something can fail. Go has the same instinct but implements it with an
ordinary return value instead of a separate control-flow path:

```java
// Java
try {
    String data = Files.readString(Path.of("file.txt"));
} catch (IOException e) {
    System.out.println(e.getMessage());
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

Go functions can return more than one value — that's what `data, err :=` is
doing. There's no `try`, no `catch`, no `finally`, and no unchecked exceptions
travelling silently up the stack. The `if err != nil` pattern appears
everywhere; it's repetitive, but every failure path is visible by reading
straight down the page.

For cleanup where you'd reach for `finally` or try-with-resources, Go has
`defer` — it schedules a call to run when the function returns:

```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close() // runs however this function exits
```

## Zero values, not null

In Java, an uninitialised object reference is `null`, and guarding against it is
a reflex. In Go, every type has a **zero value** that's ready to use:

```java
// Java
String name;         // null
int count;           // 0
List<String> items;  // null — a NullPointerException waiting to happen
```

```go
// Go
var name string    // "" — an empty string, never nil
var count int      // 0
var items []string // nil, but len(items) is 0 and append(items, "a") works
```

Go does have `nil`, but far fewer things can be it — strings, numbers, booleans
and structs all have usable zero values. Most of the null-checking habit you've
built up simply doesn't apply.

## Other differences you'll notice

| Java | Go |
|------|-----|
| `public` / `private` / `protected` | Uppercase = exported, lowercase = unexported |
| One public class per file | Many types per file; a package spans files |
| `var x = 5` or `int x = 5` | `x := 5` (short declaration) |
| Maven / Gradle, `pom.xml` | `go get` + `go.mod` |
| `System.out.println()` | `fmt.Println()` |
| `String.format()` | `fmt.Sprintf()` |
| `ArrayList<String>` | slice — `[]string` |
| `HashMap<String, String>` | map — `map[string]string` |
| `null` | `nil` (and zero values) |
| Checked exceptions, `throws` | `error` return values |
| `finally` / try-with-resources | `defer` |
| Threads, `ExecutorService`, `CompletableFuture` | Goroutines + channels (Task 04) |
| `implements` declared explicitly | Implicit — matching methods is enough |
| Checkstyle, formatter arguments | `gofmt`, built in and non-negotiable |

Concurrency is worth calling out. Where Java gives you thread pools, executors
and futures to compose, Go gives you one primitive: put `go` in front of a
function call and it runs concurrently. Goroutines are cheap enough to start
thousands of them, and they pass results to each other over **channels**. You'll
build this properly in [Task 04](../04-concurrency/).

Don't worry about memorising all of this — it'll make more sense as you work
through the exercises. The key thing: Go is deliberately simple. There's usually
one way to do something, not five.

---

Ready to start? Head to [Task 01 — Hello World](../01-hello-world/).
