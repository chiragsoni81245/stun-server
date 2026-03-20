# -------- Builder Stage --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go mod in some cases)
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Support multi-arch builds
ARG TARGETOS
ARG TARGETARCH

# Build static binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o stun-server ./cmd/server/main.go


# -------- Final Stage --------
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/stun-server .

# Create non-root user
RUN adduser -D appuser
USER appuser

# Expose ports
EXPOSE 3478/udp 3479/udp 9090

# Run binary
CMD ["./stun-server"]
