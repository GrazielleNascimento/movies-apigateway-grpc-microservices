FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create temporary build directory with write permissions
RUN mkdir -p /tmp/build && chmod 777 /tmp/build

# Set working directory
WORKDIR /tmp/build

# Copy proto files first
COPY proto /tmp/build/proto

# Copy go mod files
COPY api-gateway/go.mod api-gateway/go.sum ./

# Add replace directive for proto module
RUN go mod edit -replace=github.com/movie-microservice/proto=/tmp/build/proto

# Download dependencies
RUN go mod download

# Copy source code
COPY api-gateway/cmd ./cmd
COPY api-gateway/internal ./internal
COPY api-gateway/docs ./docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api-gateway ./cmd/main.go

# Final stage
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /tmp/build/api-gateway /api-gateway

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["/api-gateway"]
