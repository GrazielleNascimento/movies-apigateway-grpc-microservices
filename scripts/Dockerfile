FROM golang:1.22.2-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY scripts/go.mod scripts/go.sum ./

# Download dependencies
RUN go mod download


# Copy source code
COPY scripts/init_data.go .

# Copy movies.json from root to /app
COPY movies.json .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o init_data .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/init_data .
COPY movies.json .

ENTRYPOINT ["./init_data"]

