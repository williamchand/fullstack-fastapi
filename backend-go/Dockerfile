# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate code and build
RUN make generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/myapp ./cmd/server

# Runtime stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy binary
COPY --from=builder /app/bin/myapp .

# Copy migrations
COPY db/migrations ./db/migrations

# Expose ports
EXPOSE 8080 9090

# Run the application
CMD ["./myapp"]