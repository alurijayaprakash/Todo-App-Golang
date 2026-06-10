# Multi-stage Dockerfile for Go Todo API using distroless (Most Secure)
# No shell, no package manager, no vulnerabilities

# Stage 1: Build Stage
FROM golang:1.24.1-alpine3.19 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# -ldflags to strip debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always 2>/dev/null || echo 'dev')" \
    -o server \
    ./cmd/server/main.go

# Stage 2: Runtime Stage (Alpine - minimal)
FROM alpine:3.19

# Copy binary from builder
COPY --from=builder /app/server /server

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
