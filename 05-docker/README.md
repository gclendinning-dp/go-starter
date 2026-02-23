# Task 05 — Docker

## What You'll Learn

- What **containers** and **Docker** are
- The difference between an **image** and a **container**
- How to write a **Dockerfile** to package your Go server
- How to **build**, **run**, and **test** a containerised application
- What **port mapping** is and why you need it

## Key Concepts

### What is a Container?

Imagine you've baked a cake and you need to send it to a friend. You put the
cake in a sealed box with everything it needs — the plate, a fork, a napkin,
even a note with serving instructions. Your friend opens the box and has
everything required, regardless of what's already in their kitchen.

A **container** works the same way for software. It's a sealed box that contains
your program plus everything it needs to run — the right version of Go, any
libraries, configuration files, everything. When someone runs your container, it
works exactly the same on their machine as it does on yours. No more "but it
works on my computer!"

### What is Docker?

**Docker** is the tool that creates and runs containers. Think of it as:
- A **box-making machine** — it reads your instructions and builds a sealed box
- A **box-opening machine** — it can run what's inside the box

You'll use two main Docker commands:
- `docker build` — creates a container image (makes the box)
- `docker run` — starts a container from that image (opens the box and runs
  what's inside)

### Images vs Containers

This distinction trips up a lot of people, but it's simple:

- An **image** is like a **recipe card** — it describes exactly how to set up
  the environment and what to run. It's a read-only template.
- A **container** is like **the dish you cooked** using that recipe — it's a
  running instance created from the image.

You can create many containers from the same image, just like you can cook the
same recipe multiple times. Each container runs independently.

```
Image (recipe)  →  Container (running dish)
                →  Container (another copy)
                →  Container (yet another)
```

### What is a Dockerfile?

A **Dockerfile** is a text file that tells Docker how to build your image. It's
a series of instructions, each one adding a layer to the image. Here's what
each instruction does:

```dockerfile
FROM golang:1.26-alpine
```
**FROM** — the starting point. This says "start with an image that already has
Go 1.26 installed on Alpine Linux (a tiny Linux distribution)." You don't need
to install Go yourself — someone has already made an image with it.

```dockerfile
WORKDIR /app
```
**WORKDIR** — sets the working directory inside the container. All following
commands run from this directory. It's like `cd /app` but it also creates the
directory if it doesn't exist.

```dockerfile
COPY go.mod .
COPY main.go .
```
**COPY** — copies files from your computer into the container. The first
argument is the file on your machine, the second is where to put it inside the
container. The `.` means "the current working directory" (which we set to
`/app`).

```dockerfile
RUN go build -o server .
```
**RUN** — runs a command during the build. This compiles your Go code into a
binary called `server`. The `-o server` flag means "name the output file
`server`".

```dockerfile
EXPOSE 8080
```
**EXPOSE** — documents which port the container listens on. This doesn't
actually open the port — it's more like a note for anyone reading the
Dockerfile. The actual port opening happens when you run the container.

```dockerfile
CMD ["./server"]
```
**CMD** — the command that runs when the container starts. This tells Docker
"when someone runs this container, execute the `server` binary."

### Installing Docker Desktop on Mac

You need Docker Desktop installed to build and run containers. Install it with
Homebrew:

```bash
brew install --cask docker
```

After installation:
1. Open **Docker Desktop** from your Applications folder (or Spotlight search)
2. Wait for the **whale icon** to appear in your menu bar (top of the screen)
3. The whale should stop animating once Docker is ready
4. Verify it's working:

```bash
docker --version
```

You should see something like `Docker version 27.x.x`.

### What is Port Mapping?

Your containerised server listens on port 8080, but it's listening **inside the
container**. The container is like a separate building — it has its own network.
Your Mac can't reach port 8080 inside the container unless you set up **port
mapping**.

Port mapping is like **mail forwarding**. You tell Docker: "any traffic that
arrives at port 8080 on my Mac should be forwarded to port 8080 inside the
container."

```bash
docker run -p 8080:8080 my-server
```

The `-p 8080:8080` flag means `host-port:container-port`:
- Left side (`8080`) — the port on **your Mac**
- Right side (`8080`) — the port **inside the container**

They don't have to be the same number. You could use `-p 9090:8080` to map your
Mac's port 9090 to the container's port 8080. But keeping them the same is
simpler.

## Instructions

1. **Copy your server code** from Task 02.

   Copy `main.go` from `02-local-mock-server/answers/` (or your own solution)
   into this directory. Also create a `go.mod` file — see below.

2. **Create a `go.mod` file** in this directory:

   ```
   module server

   go 1.26.0
   ```

   This is needed because the Docker build happens in isolation — it doesn't
   have access to the root project's `go.mod`. The module name can be anything
   since this is a standalone binary.

3. **Write a Dockerfile** in this directory with these steps:
   - Start from `golang:1.26-alpine`
   - Set the working directory to `/app`
   - Copy `go.mod` and `main.go` into the container
   - Build the Go binary with `go build -o server .`
   - Expose port 8080
   - Set the default command to run `./server`

4. **Build and run your container:**

   ```bash
   docker build -t my-go-server .
   docker run -p 8080:8080 my-go-server
   ```

5. **Test it** in a new terminal:

   ```bash
   curl http://localhost:8080/message
   ```

   You should see: `{"message": "Hello from the backend!"}`

6. **Stop the container** with `Ctrl+C` in the terminal where it's running.

## Test It

There's no `main_test.go` for this exercise — the test is whether your
container builds, runs, and responds to HTTP requests.

An automated test script is provided in the `answers/` directory:

```bash
bash answers/test_container.sh
```

This script builds the image, starts the container, sends a request, checks
the response, and cleans everything up automatically.

## If You Get Stuck

A working Dockerfile and all supporting files are in the `answers/` directory.

## Hints

- Your Dockerfile should be exactly 7 lines (one for each instruction above)
- Make sure Docker Desktop is running (whale icon in the menu bar) before
  building
- If port 8080 is already in use, stop whatever's using it or use a different
  host port: `docker run -p 9090:8080 my-go-server` (then curl port 9090)
- To see running containers: `docker ps`
- To stop all running containers: `docker stop $(docker ps -q)`
