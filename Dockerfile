# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag for generating docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init -g cmd/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install netcat for database health check
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the wait script
COPY scripts/wait-for-db.sh /wait-for-db.sh
RUN chmod +x /wait-for-db.sh

# Expose port
EXPOSE 8080

# Use the wait script to ensure database is ready before starting
CMD ["/wait-for-db.sh", "./main"]
