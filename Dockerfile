# Stage 1: Build the Go application
# Use the official Go image based on Alpine Linux for a smaller footprint
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first
# This helps in caching dependencies separately, reducing build times
COPY go.mod go.sum ./

# Download and cache Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# Disable CGO and set the target OS and architecture for a statically linked binary
RUN export MODULE_NAME=$(go list -m) && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server $MODULE_NAME/cmd/server

# Stage 2: Create the final image
# Use a minimal base image for the final container
FROM gcr.io/distroless/base-debian11

# Set the working directory for the application
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/server .

# Copy the templates to the final image
COPY --from=builder /app/internal/features/email/templates /app/internal/features/email/templates

# Use a non-root user for security reasons
USER nonroot:nonroot

# Command to run the application
CMD ["./server"]
