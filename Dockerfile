# Builder stage
FROM --platform=$BUILDPLATFORM golang:1.23-alpine

# Install build dependencies
RUN apk add --no-cache wget tar

# Install Redis
RUN apk add --no-cache redis

# Download and install IPFS for linux-amd64
RUN apk add kubo

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Expose ports
EXPOSE 3000
EXPOSE 5001
EXPOSE 4001
EXPOSE 8080
EXPOSE 6379

# Start Redis and IPFS, then your application
CMD redis-server --bind 0.0.0.0 --port 6379 --dir /tmp & ipfs init --profile server && ipfs daemon --enable-pubsub-experiment & sleep 5 && ./main
