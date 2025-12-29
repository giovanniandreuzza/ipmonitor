# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o ipmonitor \
    ./cmd/ipmonitor

# Final stage
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /build/ipmonitor /ipmonitor

# Create directory for IP storage (will be mounted as volume)
VOLUME ["/data"]

# Set environment variable for IP storage location
ENV IP_STORAGE_PATH=/data/public_ip.txt

ENTRYPOINT ["/ipmonitor"]
