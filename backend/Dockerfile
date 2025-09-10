# ----------------------------
# Builder Stage
# ----------------------------
FROM golang:1.23-alpine AS builder

# Set Go environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o /erebus ./cmd/erebusd

# ----------------------------
# Runtime Stage
# ----------------------------
FROM alpine:3.18

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder stage
COPY --from=builder /erebus /usr/local/bin/erebus

# Expose application port
EXPOSE 8080

# Set default entrypoint
ENTRYPOINT ["/usr/local/bin/erebus"]
