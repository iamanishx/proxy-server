# proxy-server

> This server can proxy upto 100k requests per sec 🚀

A simple, high-performance HTTP reverse proxy server written in Go. It forwards incoming HTTP requests to a configurable backend server with support for connection pooling and idle connection management.

## Prerequisites

- Go 1.23.3 or later

## Configuration

The proxy server is configured via environment variables loaded from a `.env` file:

| Variable                  | Default         | Description                                  |
| ------------------------- | --------------- | -------------------------------------------- |
| `LISTEN_ADDR`             | `:8080`         | Address and port the proxy listens on        |
| `BACKEND_URL`             | _(required)_    | URL of the backend service to proxy to       |
| `MAX_IDLE_CONNS`          | `1000`          | Maximum number of idle connections in the pool |
| `MAX_IDLE_CONNS_PER_HOST` | `1000`          | Maximum idle connections per backend host    |

## Quick Start

1. **Clone the repository:**
   ```bash
   git clone <repo-url>
   cd proxy-server
   ```

2. **Configure the `.env` file:**
   ```env
   LISTEN_ADDR=:8080
   BACKEND_URL=http://localhost:8080
   MAX_IDLE_CONNS=1000
   MAX_IDLE_CONNS_PER_HOST=1000
   ```

3. **Run the server:**
   ```bash
   go run ./cmd/main.go
   ```

4. **Send a request:**
   ```bash
   curl http://localhost:8080/any/path
   ```

## Makefile Targets

| Target    | Description                              |
| --------- | ---------------------------------------- |
| `fmt`     | Run `go fmt ./...`                       |
| `vet`     | Run `go vet ./...`                       |
| `tidy`    | Run `go mod tidy`                        |
| `lint`    | Run `golangci-lint run ./...`            |
| `build`   | Build the binary to `bin/proxy-server`   |
| `run`     | Run the server with `go run`             |
| `test`    | Run `go test ./...`                      |
| `check`   | Run fmt, vet, tidy, and lint in sequence |
| `clean`   | Remove the built binary                  |
