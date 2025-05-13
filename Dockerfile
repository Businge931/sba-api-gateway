FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/main.go

# Use a minimal alpine image for the final stage
FROM alpine:3.17

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /api-gateway .

# Expose the port the service runs on
EXPOSE 8080

# Run the service
CMD ["./api-gateway"]
