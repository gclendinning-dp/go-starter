# Task 09 — Docker Compose: Replicas & Load Balancing

## What You'll Learn

- What **Docker Compose** is and how it manages multiple containers
- What **replicas** are and why you'd run multiple copies of a service
- What a **load balancer** is and how Nginx distributes traffic
- How `os.Hostname()` works inside containers
- The difference between `expose` and `ports` in Docker Compose

## Key Concepts

### What is Docker Compose?

In Task 05 you ran a single container — like a single food truck serving
customers. That works great for one service, but real applications usually need
multiple services working together (a web server, a database, a cache, etc.).

**Docker Compose** is like managing a **food court** instead of a single food
truck. You write one file (`docker-compose.yml`) that describes all the stalls,
and Compose sets them all up, connects them together, and tears them all down
with a single command.

Instead of running several `docker build` and `docker run` commands, you just:

```bash
docker compose up -d --build    # start everything
docker compose down             # stop everything
```

### What are Replicas?

Imagine a busy coffee shop with one cashier. There's a queue. Now add two more
cashiers doing the exact same job — the queue moves three times faster.

**Replicas** are the same idea. Instead of running one copy of your server, you
run three (or more). Each replica is identical — same code, same image — but
they run in separate containers. Docker Compose creates them for you:

```yaml
deploy:
  replicas: 3
```

### What is a Load Balancer?

If you have three cashiers, someone needs to direct customers to a free one.
That's the **load balancer** — a receptionist who directs incoming visitors to
an available office.

In this exercise, **Nginx** is the load balancer. It sits in front of your three
server replicas and distributes incoming requests across them. By default, Nginx
uses **round-robin** — it sends the first request to server 1, the second to
server 2, the third to server 3, then back to server 1, and so on.

### How Nginx Load Balancing Works

The Nginx configuration uses an **upstream** block:

```nginx
upstream backend {
    server web:8080;
}
```

Wait — there's only one `server web:8080` line, but we have three replicas?
Here's the clever bit: Docker Compose gives all replicas of the `web` service
the same DNS name (`web`). When Nginx resolves `web`, Docker's DNS returns
**all three container IPs**, and Nginx automatically distributes traffic across
them.

### `os.Hostname()` in Containers

When you run `os.Hostname()` on your Mac, it returns your computer's name. But
inside a container, it returns the **container ID** — a short hex string like
`a1b2c3d4e5f6`. Each replica gets a different container ID, so this is a handy
way to prove that different containers are handling your requests.

### `expose` vs `ports`

In Docker Compose, there are two ways to make a port available:

- **`expose`** — makes the port available **only to other containers** in the
  same Compose network. Outside world (your Mac) can't reach it directly.
- **`ports`** — maps the container's port to your Mac, so you can access it
  from `localhost`.

In this exercise, the `web` replicas **expose** port 8080 (only Nginx can reach
them), and Nginx **maps** port 80 to your Mac's port 8080 (so you can reach
Nginx from `localhost:8080`).

## Instructions

1. **Modify the Task 07 server** to include the hostname in responses.

   Create a `main.go` in this directory (not in `answers/`) based on your Task 07
   server with these changes:

   - `MessageHandler` returns `{"message": "Hello from <hostname>!", "hostname": "<hostname>"}`
     using `os.Hostname()` and `json.NewEncoder(w).Encode()`
   - `GreetHandler` returns `{"greeting": "Hello, Alice! From <hostname>", "hostname": "<hostname>"}`

   You'll need to define response structs with a `Hostname` field.

2. **Create a `go.mod`** in this directory:

   ```
   module server

   go 1.26.0
   ```

3. **Create a `Dockerfile`** — same pattern as Task 05:

   ```dockerfile
   FROM golang:1.26-alpine
   WORKDIR /app
   COPY go.mod .
   COPY main.go .
   RUN go build -o server .
   EXPOSE 8080
   CMD ["./server"]
   ```

4. **Create an `nginx.conf`** file:

   ```nginx
   upstream backend {
       server web:8080;
   }

   server {
       listen 80;

       location / {
           proxy_pass http://backend;
       }
   }
   ```

   This tells Nginx: "for any request, pass it to one of the `web` containers
   on port 8080."

5. **Create a `docker-compose.yml`**:

   ```yaml
   services:
     web:
       build: .
       deploy:
         replicas: 3
       expose:
         - "8080"

     nginx:
       image: nginx:alpine
       ports:
         - "8080:80"
       volumes:
         - ./nginx.conf:/etc/nginx/conf.d/default.conf
       depends_on:
         - web
   ```

   Key points:
   - `build: .` tells Compose to build the image from the Dockerfile in this
     directory
   - `replicas: 3` creates three containers from the same image
   - `expose` makes port 8080 available only inside the Compose network
   - `volumes` mounts your nginx.conf into the Nginx container
   - `depends_on` ensures `web` starts before `nginx`

6. **Start the stack:**

   ```bash
   docker compose up -d --build
   ```

7. **Test load balancing** — hit the endpoint multiple times:

   ```bash
   for i in $(seq 1 6); do
       curl -s http://localhost:8080/message | python3 -c "import sys,json; print(json.load(sys.stdin)['hostname'])"
   done
   ```

   You should see different hostnames, proving different containers handle
   different requests.

8. **Stop everything:**

   ```bash
   docker compose down --rmi local
   ```

## Test It

An automated test script is provided in the `answers/` directory:

```bash
bash answers/test_compose.sh
```

This script builds the stack, verifies that multiple hostnames appear (proving
load balancing), tests the POST endpoint through the load balancer, and cleans
everything up.

## If You Get Stuck

A working solution is in the `answers/` directory.

## Hints

- Get the hostname: `hostname, _ := os.Hostname()`
- Encode JSON to the response: `json.NewEncoder(w).Encode(resp)`
- Format with hostname: `fmt.Sprintf("Hello from %s!", hostname)`
- Start Compose: `docker compose up -d --build`
- Stop Compose: `docker compose down --rmi local`
- See running containers: `docker compose ps`
- See logs: `docker compose logs web`
- The imports you'll need: `"encoding/json"`, `"fmt"`, `"net/http"`, `"os"`
