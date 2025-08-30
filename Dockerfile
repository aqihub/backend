# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    redis \
    kubo

# Copy the binary from builder stage
COPY --from=builder /app/main /app/main

# Set working directory
WORKDIR /app


# Expose ports
EXPOSE 3000 5001 4001 8080 6379

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Use a startup script to handle initialization
CMD redis-server --bind 0.0.0.0 --port 6379 --dir /tmp & ipfs init --profile server && ipfs daemon --enable-pubsub-experiment & sleep 5 && ./main