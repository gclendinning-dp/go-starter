# Task 06 — File I/O

## What You'll Learn

- How to **read** a file line by line with `os.Open` and `bufio.Scanner`
- How to **write** lines to a file with `os.Create` and `fmt.Fprintln`
- How to **filter** text with `strings.Contains`
- Why `defer` matters for closing files

## Key Concepts

### What is a File?

A **file** is just data stored on your computer's disk. So far, all the data in
your programs has lived in memory — it disappears when the program exits. Files
are how you make data **persist**. Log files, configuration files, databases —
they're all files underneath.

In this exercise, you'll read a server log file, filter out the error lines,
and write them to a new file. This is the kind of thing real developers do all
the time — processing logs, transforming data, generating reports.

### Reading a File: `os.Open` + `bufio.Scanner`

To read a file in Go, you **open** it and then **scan** it line by line:

```go
file, err := os.Open("server.log")
if err != nil {
    return nil, err
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}
```

Think of `os.Open` as opening a book. `bufio.Scanner` is like putting your
finger on the first line and moving it down one line at a time. Each call to
`scanner.Scan()` moves to the next line, and `scanner.Text()` gives you the
text on that line.

### Writing a File: `os.Create` + `fmt.Fprintln`

To write lines to a file, you **create** it and then **print** into it:

```go
file, err := os.Create("errors.txt")
if err != nil {
    return err
}
defer file.Close()

for _, line := range lines {
    fmt.Fprintln(file, line)
}
```

`os.Create` is like getting a fresh blank page. If a file with that name
already exists, it gets replaced. `fmt.Fprintln` works exactly like
`fmt.Println` except it writes to the file instead of the screen — the `F`
stands for "file" (well, technically any `io.Writer`, but think "file" for now).

### Filtering with `strings.Contains`

`strings.Contains` checks if one string appears inside another:

```go
strings.Contains("ERROR: disk full", "ERROR")  // true
strings.Contains("INFO: all good", "ERROR")    // false
```

It's like using Ctrl+F (or Cmd+F) in a document — you give it a keyword and
it tells you whether it found a match.

### Why `defer file.Close()`?

When you open a file, the operating system allocates resources to track it.
If you forget to close the file, those resources leak. In a long-running
program, this can eventually crash your application.

`defer file.Close()` means "close this file when the function returns, no
matter what". You put it right after opening the file so you never forget:

```go
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close() // runs when the function exits
// ... use the file ...
```

You already saw `defer` in Task 03 with `defer resp.Body.Close()` — same idea,
different resource.

## Instructions

1. Create a file called `main.go` in this directory (not in `answers/`).

2. Write a function called `ReadLines` that:
   - Takes a file path (`string`)
   - Returns `([]string, error)`
   - Opens the file and reads it line by line using `bufio.Scanner`
   - Returns only **non-empty** lines (skip blank lines)

3. Write a function called `FilterLines` that:
   - Takes a slice of lines (`[]string`) and a keyword (`string`)
   - Returns `[]string`
   - Returns only the lines that **contain** the keyword

4. Write a function called `WriteLines` that:
   - Takes a file path (`string`) and a slice of lines (`[]string`)
   - Returns `error`
   - Creates the file and writes each line into it

5. In your `main` function:
   - Read `server.log` using `ReadLines`
   - Filter for lines containing `"ERROR"` using `FilterLines`
   - Write the results to `errors.txt` using `WriteLines`
   - Print how many errors were found

## Run It

```bash
go run .
```

You should see something like: `Found 3 errors, written to errors.txt`

Then check the output:

```bash
cat errors.txt
```

## Test It

Copy `main_test.go` from the `answers/` directory into this directory, then:

```bash
go test -v
```

The tests check four things:
1. **ReadLines** — reads a file and returns the correct lines
2. **FilterLines** — filters lines by keyword correctly
3. **WriteLines** — writes lines to a file that can be read back
4. **ReadLines (not found)** — returns an error for a nonexistent file

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Open a file: `os.Open(path)`
- Create a file: `os.Create(path)`
- Scan lines: `bufio.NewScanner(file)`, then `scanner.Scan()` and `scanner.Text()`
- Check if a string contains a keyword: `strings.Contains(line, keyword)`
- Write a line to a file: `fmt.Fprintln(file, line)`
- Don't forget `defer file.Close()` after opening or creating a file
- The imports you'll need: `"bufio"`, `"fmt"`, `"os"`, `"strings"`
