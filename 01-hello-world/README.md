# Task 01 — Hello World

## What You'll Learn

- Every Go program starts with `package main`
- How to import packages from the standard library
- How to use `fmt.Println` to print output
- What **strings** and **functions** are

## Key Concepts

### What is a String?

A **string** is just text. In Go, you write strings inside double quotes:

```go
"Hello, world!"
"Alice"
"This is a string too."
```

Strings can contain letters, numbers, spaces, punctuation — anything you can
type. They're one of the most basic building blocks in programming.

### What is a Function?

A **function** is a reusable block of code that does a specific job. Think of it
like a recipe: you give it ingredients (called **parameters**), it does some
work, and it gives you back a result (called a **return value**).

```go
func Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

Breaking this down:

- `func` — tells Go you're defining a function
- `Greet` — the name of the function (you choose this)
- `(name string)` — the input: a parameter called `name` that must be a string
- `string` (after the parentheses) — the output: this function returns a string
- `return` — sends a value back to whoever called the function

### What is `fmt.Sprintf`?

`fmt.Sprintf` builds a string by filling in placeholders. The `%s` placeholder
gets replaced with a string value:

```go
fmt.Sprintf("Hello, %s! Welcome to Go.", "Alice")
// Result: "Hello, Alice! Welcome to Go."
```

It's like a template with blanks to fill in.

### What is a Package?

Go organises code into **packages**. Every `.go` file starts with a package
declaration. The special package `main` tells Go "this is a program you can
run" (as opposed to a library that other code uses).

`import "fmt"` brings in the `fmt` package from Go's standard library, which
gives you printing and formatting tools like `Println` and `Sprintf`.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).
2. Your file needs three things:

   ```go
   package main    // declares this is a runnable program

   import "fmt"    // imports the formatting/printing package

   func main() {   // the entry point — Go runs this function first
       // your code here
   }
   ```

3. Write a function called `Greet` that takes a `name` (a string) and returns
   a greeting string in this format: `"Hello, <name>! Welcome to Go."`

   For example, `Greet("Alice")` should return `"Hello, Alice! Welcome to Go."`

4. In your `main` function, call `Greet` with your own name and print the result.

## Run It

```bash
go run .
```

## Test It

Copy `main_test.go` from the `answers/` directory into this directory, then:

```bash
go test -v
```

The test will check that your `Greet` function returns the correct string for
several different inputs.

## If You Get Stuck

A working solution is in the `answers/` directory. Try to solve it yourself
first — you'll learn more that way — but it's there if you need it.

## Hints

- Use `fmt.Sprintf` to build a string with variables: `fmt.Sprintf("Hello, %s!", name)`
- The function signature should look like: `func Greet(name string) string`
