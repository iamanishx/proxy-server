# proxy-server

> A high-performance HTTP reverse proxy server written in Go.

## Features

- Reverse proxy forwarding to a configurable backend URL
- Configurable maximum idle connections and idle connections per host
- Environment-based configuration via `.env` file
- Connection pooling with configurable timeouts
- Lightweight and efficient

## Prerequisites

- Go 1.23.3 or later

## Installation

Clone the repository and build the binary:

```bash
git clone <repo-url> && cd proxy-server
make build
```

Or build directly with Go:

```bash
go build -o bin/proxy-server ./cmd/main.go
```

## Usage

1. Create a `.env` file in the project root (or copy the provided `.env`).
2. Set the required environment variables.
3. Run the server:

```bash
make run
```

Or directly:

```bash
go run ./cmd/main.go
```

## Configuration

The following environment variables are read from the `.env` file:

| Variable                  | Default  | Description                                 |
|---------------------------|----------|---------------------------------------------|
| `BACKEND_URL`             | —        | **Required.** The backend URL to proxy to   |
| `LISTEN_ADDR`             | `:8080`  | The address the proxy server listens on     |
| `MAX_IDLE_CONNS`          | `1000`   | Maximum number of idle connections          |
| `MAX_IDLE_CONNS_PER_HOST` | `1000`   | Maximum idle connections per host           |

Example `.env`:

```
BACKEND_URL=http://localhost:3000
LISTEN_ADDR=:8080
MAX_IDLE_CONNS=1000
MAX_IDLE_CONNS_PER_HOST=1000
```

## Makefile Targets

| Target    | Description                                |
|-----------|--------------------------------------------|
| `fmt`     | Format Go source code                      |
| `vet`     | Run Go vet                                 |
| `tidy`    | Tidy Go module dependencies                |
| `lint`    | Run golangci-lint                          |
| `build`   | Build the binary to `bin/proxy-server`     |
| `run`     | Run the server (go run)                    |
| `test`    | Run all tests                              |
| `check`   | Run fmt, vet, tidy, and lint sequentially  |
| `clean`   | Remove built binaries                      |

## License

MIT
