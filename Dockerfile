# Stage 1: Build the Go application
# We use a Go-specific image to compile the source code.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application binary
# The `-ldflags "-s -w"` flag strips debug information for a smaller binary size.
# The `-a -installsuffix cgo` flag helps create a statically-linked binary,
# making it portable and suitable for a minimal base image.
RUN go build -o main ./cmd/server/main.go

# Stage 2: Create a lightweight image for the final application
# We use a minimal Alpine image as the base for the final container.
FROM gcr.io/distroless/base-debian11

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the Go application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
