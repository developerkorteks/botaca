# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bot cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite tzdata

# Set timezone
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bot .

# Create directories for data persistence
RUN mkdir -p /app/data

# Set permissions
RUN chmod +x ./bot

# Expose port (if needed for health checks)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD pgrep -f "./bot" || exit 1

# Run the application
CMD ["./bot"]