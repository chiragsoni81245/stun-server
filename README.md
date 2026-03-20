# STUN Server (Go)

A production-ready STUN (Session Traversal Utilities for NAT) server written in Go. This server helps clients discover their public IP address and port, enabling NAT traversal and network diagnostics.

---

## Features

* STUN Binding Request support (RFC 5389 compliant via pion/stun)
* Multi-socket support (multiple ports)
* Per-IP rate limiting
* Prometheus metrics endpoint
* Structured logging (zerolog)
* Graceful shutdown
* Configurable via flags

---

## Project Structure

```id="pzctri"
stun-server/
├── cmd/server/main.go
├── internal/stunserver/
│   ├── server.go
│   ├── ratelimiter.go
│   ├── metrics.go
│   └── logger.go
├── pkg/config/config.go
├── Dockerfile
└── go.mod
```

---

## Getting Started

### Prerequisites

* Go 1.22+
* A machine with a public IP (for real-world usage)

---

### Installation

1. Use pre build binary from release
https://github.com/chiragsoni81245/stun-server/releases

2. Build from source
```id="93oexn"
git clone https://github.com/chiragsoni81245/stun-server.git
cd stun-server
go mod tidy
go build -o bin/stun-server ./cmd/server/main.go
```

---

### Run Locally

```id="faz7js"
./stun-server ./cmd/server
```

---

### Run with Custom Ports

```id="drs9pe"
./stun-server --addrs=":3478,:3479"
```

---

## Docker

### Run

```id="ywe3wm"
docker run -p 3478:3478/udp -p 3479:3479/udp -p 9090:9090 chiragsoni812/stun-server
```

---

## Configuration

| Flag    | Description                          | Default       |
| ------- | ------------------------------------ | ------------- |
| `addrs` | Comma-separated UDP listen addresses | `:3478,:3479` |

---

## Metrics

Metrics are exposed at:

```id="9uf8f1"
http://localhost:9090/metrics
```

### Available Metrics

* `stun_requests_total` – Total STUN requests
* `stun_errors_total` – Total errors
* `stun_request_duration_seconds` – Request latency histogram

---

## Example Prometheus Configuration

```id="tlef17"
scrape_configs:
  - job_name: "stun-server"
    static_configs:
      - targets: ["localhost:9090"]
```

---

## Testing

You can test the server using any STUN client or your own NAT testing tool:

```id="mrfzy7"
./net-tester --stun <server-ip>:3478
```

---

## Notes

* This server implements basic STUN functionality (Binding Requests only).
* It does not support TURN (relay functionality).
* For accurate NAT type detection, multiple STUN endpoints (ports or IPs) are recommended.

---
