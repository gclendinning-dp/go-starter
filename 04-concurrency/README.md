# Task 04 — Concurrency

## What You'll Learn

- What **concurrency** means and why it's useful
- How to run code simultaneously with **goroutines**
- How to collect results using **channels**
- How to wait for everything to finish with **sync.WaitGroup**
- What **slices** are (Go's expandable lists)

## Key Concepts

### What is a Slice?

Before we get to concurrency, you need to know about **slices**. A slice is
Go's version of an expandable list. You've used single values like `string` and
`int` — a slice holds multiple values of the same type.

```go
names := []string{"Alice", "Bob", "Charlie"}
```

The `[]string` part means "a slice of strings". You can access items by their
position (starting from 0):

```go
fmt.Println(names[0]) // prints: Alice
fmt.Println(names[1]) // prints: Bob
```

You can add items with `append`:

```go
names = append(names, "Diana")
// names is now ["Alice", "Bob", "Charlie", "Diana"]
```

And you can loop over a slice with `range`:

```go
for i, name := range names {
    fmt.Println(i, name)
}
// 0 Alice
// 1 Bob
// 2 Charlie
// 3 Diana
```

You can also check how many items are in a slice with `len`:

```go
fmt.Println(len(names)) // prints: 4
```

### What is `time.Sleep`?

`time.Sleep` pauses your program for a specified duration. It's useful for
simulating slow operations like network requests:

```go
time.Sleep(200 * time.Millisecond) // pause for 200ms
```

In real code you'd make an actual HTTP request, but for learning concurrency
it's cleaner to simulate the delay. This lets you focus on the concurrency
concepts without worrying about servers and URLs.

### What is Concurrency?

Imagine you're a chef in a kitchen. You need to:
1. Boil pasta (10 minutes)
2. Make sauce (8 minutes)
3. Prepare salad (5 minutes)

If you do these **one at a time** (sequentially), it takes 23 minutes. But
you're smarter than that — you put the pasta on to boil, then start the sauce
while it's cooking, and prepare the salad while both are going. The total time
drops to about 10 minutes because you're doing things **concurrently**.

That's concurrency: managing multiple tasks that are **in progress at the same
time**. Your program doesn't wait for one thing to finish before starting the
next.

### What is a Goroutine?

A **goroutine** is Go's way of running a function concurrently. It's
incredibly simple — you just put the word `go` before a function call:

```go
go doSomething()
```

That's it. The function `doSomething` now runs in the background while your
program continues to the next line.

Think of it like asking a friend to go grab you a coffee. You say "go get me a
coffee" and then you carry on with your work. You don't stand there watching
them walk to the shop — you do other things and they'll come back with the
coffee eventually.

Without `go`, you'd be walking to the shop yourself — everything else stops
until you get back.

```go
// Sequential — one after another (slow)
result1 := SlowFetch("url1") // wait 200ms...
result2 := SlowFetch("url2") // wait another 200ms...
// Total: 400ms

// Concurrent — both at the same time (fast)
go SlowFetch("url1") // starts immediately
go SlowFetch("url2") // also starts immediately
// Total: ~200ms (they run in parallel)
```

### What is a Channel?

When a goroutine finishes its work, how do you get the result back? You use a
**channel**.

Think of a channel like a **letterbox**. One goroutine puts a letter in (sends
a value), and another goroutine takes the letter out (receives the value):

```go
ch := make(chan string) // create a channel that carries strings

// Goroutine sends a value into the channel
go func() {
    ch <- "hello" // put "hello" into the letterbox
}()

// Main function receives the value
msg := <-ch // take the letter out of the letterbox
fmt.Println(msg) // prints: hello
```

The `<-` arrow shows the direction:
- `ch <- value` means "send value into the channel" (put letter in the box)
- `value := <-ch` means "receive a value from the channel" (take letter out)

Receiving from a channel **blocks** — your program pauses and waits until
something arrives. This is actually useful because it naturally synchronises
your code.

### What is a WaitGroup?

A **WaitGroup** is like a teacher counting children on a school field trip.
Before setting off, the teacher counts how many children there are. As each
child gets back on the bus, the count goes down. The teacher waits until the
count reaches zero before driving off.

```go
var wg sync.WaitGroup

wg.Add(3)    // "we're waiting for 3 things"
// ... start 3 goroutines ...
// Each goroutine calls wg.Done() when it finishes ("I'm back on the bus!")
wg.Wait()    // blocks until the count reaches zero ("is everyone here?")
```

In this exercise, you'll use a WaitGroup to know when all your goroutines have
finished sending their results into the channel.

### Putting It Together

Here's the pattern you'll use:

1. Create a channel to collect results
2. Create a WaitGroup to track goroutines
3. For each URL, start a goroutine that:
   - Calls `SlowFetch` to get the result
   - Sends the result into the channel
   - Tells the WaitGroup it's done
4. Start a separate goroutine that waits for all others to finish, then
   closes the channel
5. Collect all results from the channel into a slice

Closing the channel (step 4) is important — it tells the receiving side "no
more values are coming". Without it, your program would wait forever for more
results.

## Before You Start — Read the Test

Open `main_test.go` in this directory. It shows you what `FetchAll` needs to do,
how correctness is checked, and the concurrency timing assertion — 5 tasks at
200ms each must complete in under 500ms.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. The `SlowFetch` function is provided for you — copy it exactly:

   ```go
   // SlowFetch simulates a slow network request. It takes 200ms and returns
   // a string based on the URL.
   func SlowFetch(url string) string {
       time.Sleep(200 * time.Millisecond)
       return "response from " + url
   }
   ```

3. Write a function called `FetchAll` that:
   - Takes a `urls` parameter of type `[]string` (a slice of strings)
   - Returns a `[]string` (a slice of results)
   - Runs `SlowFetch` concurrently for **each** URL using goroutines
   - Collects all results via a channel
   - Uses a `sync.WaitGroup` to wait for all goroutines to finish
   - Returns the collected results (the order doesn't matter)

4. In your `main` function, call `FetchAll` with a few example URLs and print
   the results.

## Run It

```bash
go run .
```

## Test It

```bash
go test -v
```

The tests check three things:
1. **Correctness** — each URL produces the expected result
2. **Concurrency** — 5 URLs (each taking 200ms) complete in under 500ms,
   proving they ran concurrently rather than one-at-a-time
3. **Empty input** — passing an empty slice returns an empty slice without
   errors

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Create a channel: `ch := make(chan string)`
- Start a goroutine: `go func() { ... }()`
- Send to a channel: `ch <- value`
- Receive from a channel in a loop: `for result := range ch { ... }`
- Close a channel: `close(ch)`
- Don't forget to call `wg.Add(1)` before each goroutine and `wg.Done()`
  inside each goroutine (use `defer wg.Done()` at the top)
- Close the channel after all goroutines finish:
  ```go
  go func() {
      wg.Wait()
      close(ch)
  }()
  ```
- The imports you'll need: `"sync"`, `"time"`
