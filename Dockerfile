# --------------------------------------------------------
# Stage 1: Build the Application
# --------------------------------------------------------
FROM golang:alpine AS builder

WORKDIR /app

# Copy the dependency files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary named "limiter"
# CGO_ENABLED=0 ensures a static binary (no external C libraries needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o limiter ./cmd/api

# --------------------------------------------------------
# Stage 2: The Final Tiny Image
# --------------------------------------------------------
FROM alpine:latest

WORKDIR /root/

# Copy ONLY the binary from the builder stage
COPY --from=builder /app/limiter .

# Expose the port
EXPOSE 8080

# Command to run
CMD ["./limiter"]