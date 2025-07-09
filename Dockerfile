# Stage 1: Build the Go app
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (used for fetching dependencies), add ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the source code
COPY . .

# Build the Go binary
RUN go build -o server .


# Stage 2: Run the app
FROM alpine:latest

WORKDIR /root/

# Copy the compiled Go binary from builder
COPY --from=builder /app/server .

# Swagger files (optional but good)
COPY --from=builder /app/docs ./docs

# Expose the Gofr port (default 9000) and Swagger UI (on same port)
EXPOSE 9000

CMD ["./server"]
