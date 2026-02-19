# Distributed Rate Limiter (Go + Redis + Lua)

A high-performance, thread-safe rate limiter middleware for Go APIs. It implements the **Token Bucket** algorithm using **Redis Lua Scripts** to ensure atomicity in a distributed environment.

This project demonstrates how to handle concurrency, race conditions, and service discovery in a microservices architecture.

## 🚀 Key Features

* **Distributed Design:** Uses Redis as a centralized state store, allowing multiple API instances to share the same rate limits.
* **Atomic Operations:** Uses Lua scripting to prevent "Race Conditions" (Time-of-Check vs. Time-of-Use) that occur with standard `GET`/`SET` operations.
* **Token Bucket Algorithm:** Implements "Lazy Refill" logic for smooth traffic shaping (handling bursts better than Fixed Window counters).
* **Middleware Pattern:** Plugs easily into any standard `http.Handler`.
* **Containerized:** Includes a Multi-Stage Dockerfile for a lightweight production image (<15MB).

## 🛠️ Tech Stack

* **Language:** Go (Golang) 1.22+
* **Database:** Redis (State Management)
* **Scripting:** Lua (Atomic Transactions)
* **Infrastructure:** Docker & Docker Networks
* **Library:** `go-redis/v9`

## 📂 Project Structure

```text
rate-limiter/
├── cmd/
│   └── api/
│       └── main.go           # Application Entry Point
├── internal/
│   ├── limiter/
│   │   └── limiter.go        # Redis/Lua Token Bucket Logic
│   └── middleware/
│       └── rate_limit.go     # HTTP Middleware Wrapper
├── Dockerfile                # Multi-stage build definition
└── go.mod
```

# ⚙️ How It Works (The Algorithm)

Instead of a simple counter, we simulate a bucket of tokens:

* **Capacity:** 5 tokens (Max Burst).
* **Refill Rate:** 1 token every 12 seconds (approx 5/min).

### 1. Lazy Refill Strategy
We do **not** use a background timer to add tokens. Instead, we calculate the refill amount only when a request actually arrives. 

When a user hits the API, the system calculates:
1.  `delta = now - last_request_time`
2.  `new_tokens = old_tokens + (delta * refill_rate)`

### 2. Atomicity (The "Race Condition" Fix)
All of this math happens inside a single **Lua Script** running directly on Redis. 

This is critical because standard `GET` and `SET` operations are not atomic. Without Lua, two requests hitting the API at the exact same microsecond could both read "4 tokens" and incorrectly decrement it to "3", allowing more traffic than intended. The Lua script guarantees that no other command runs while the balance is being updated.

# 🏃‍♂️ How to Run

### Option 1: Docker (Recommended)

Since this project requires Redis and the Go application to communicate, we use a Docker Bridge Network.

1.  **Create a Bridge Network:**
    ```bash
    docker network create app-net
    ```

2.  **Start Redis:**
    ```bash
    docker run -d --name my-redis --network app-net redis
    ```

3.  **Build the Rate Limiter Image:**
    ```bash
    docker build -t rate-limiter .
    ```

4.  **Run the Rate Limiter:**
    *We pass the Redis address using an environment variable so the app knows where to find the database.*
    ```bash
    docker run -p 8080:8080 \
      --network app-net \
      -e REDIS_ADDR="my-redis:6379" \
      rate-limiter
    ```

---

### Option 2: Local Development

If you have Go and Redis installed on your machine:

1.  **Start Redis:**
    Ensure Redis is running on the default port (`localhost:6379`).

# 🧪 API Usage

Once the server is running (on port 8080), you can test the rate limiter using `curl` or a browser.

**Test Endpoint:** `GET /?user={id}`

1. Success Response (200 OK)
2. Blocked Response (429 Too Many Requests)
