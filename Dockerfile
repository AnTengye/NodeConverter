# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/nodeconverter .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/nodeconverter /app/nodeconverter

# Copy config files
COPY config.yaml /app/config.yaml
COPY clash-tpl.yaml /app/clash-tpl.yaml

# Expose the port
EXPOSE 25500

# Run the application
ENTRYPOINT ["/app/nodeconverter"]
CMD ["-f", "/app/config.yaml"]
