FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app


# Copy proto files first
COPY proto /app/proto

# Copy go mod files
COPY movies-service/go.mod movies-service/go.sum ./

# Add replace directive for proto module
RUN go mod edit -replace=github.com/movie-microservice/proto=/app/proto

# Download dependencies
RUN go mod download

# Copy source code
COPY movies-service/cmd ./cmd
COPY movies-service/internal ./internal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o movies-service ./cmd/main.go

# Final stage
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/movies-service /movies-service

# Expose gRPC port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/movies-service", "health"] || exit 1

# Run the binary
ENTRYPOINT ["/movies-service"]