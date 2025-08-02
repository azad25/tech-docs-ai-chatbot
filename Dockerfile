# Multi-stage Dockerfile for Tech Docs AI

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker cmd/worker/main.go

# Development stage
FROM golang:1.23-alpine AS development

# Install development tools
RUN apk add --no-cache git ca-certificates tzdata curl

# Install air for hot reloading (compatible version)
RUN go install github.com/cosmtrek/air@v1.49.0

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose ports
EXPOSE 8080

# Default command for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:latest AS production

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy built binaries from builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/worker .

# Copy configuration files if any
COPY --from=builder /app/nginx.conf ./nginx.conf

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command (can be overridden)
CMD ["./server"]