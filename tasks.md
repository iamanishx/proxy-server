# Integration Tests Implementation Plan

## Project Overview

This is a Go HTTP reverse proxy server that forwards incoming requests to a configurable backend with connection pooling. Currently, there are **no tests**. The goal is to write integration tests that verify the proxy behavior end-to-end.

## Checklist

### 1. Project Setup & Dependencies
- [x] Create an `internal/proxy/` package (or use `pkg/proxy/`) to refactor the proxy logic from `cmd/main.go` into reusable code
- [x] Add `github.com/stretchr/testify` as a test dependency for assertions (`go get github.com/stretchr/testify`)
- [x] Create a `tests/` or `integration/` directory for integration test files
- [x] Ensure the test file uses the `_test.go` suffix and has `package main` or appropriate test package

### 2. Integration Test: Basic Request Forwarding
- [x] Start a local test backend (e.g., `httptest.NewServer`) that returns known responses
- [x] Start the proxy server pointing to the test backend
- [x] Send an HTTP GET request through the proxy and verify the response status code, body, and headers match the backend's response
- [x] Send an HTTP POST with a body and verify it is forwarded correctly

### 3. Integration Test: Path Preservation
- [x] Start a test backend that echoes the request path
- [x] Send a request with a specific path (e.g., `/api/v1/resource`) through the proxy
- [x] Verify the backend receives the exact same path (i.e., the proxy does not modify the path)

### 4. Integration Test: Query Parameters Forwarding
- [x] Start a test backend that echoes query parameters
- [x] Send a request with query parameters (e.g., `?key=value&page=1`) through the proxy
- [x] Verify the backend receives the exact same query parameters

### 5. Integration Test: Header Forwarding
- [x] Start a test backend that echoes all received headers
- [x] Send a request with custom headers through the proxy
- [x] Verify the backend receives the expected headers (e.g., `Content-Type`, `X-Custom-Header`)
- [x] Verify Hop-by-hop headers (e.g., `Connection`, `Transfer-Encoding`) are stripped as expected

### 6. Integration Test: Backend Unavailable / Error Handling
- [x] Point the proxy to a backend that is down or refuses connections
- [x] Send a request through the proxy
- [x] Verify the proxy returns a `502 Bad Gateway` or appropriate error status

### 7. Integration Test: Concurrent Requests (Performance/Stability)
- [x] Start a test backend with a small delay on each request
- [x] Send multiple concurrent requests through the proxy (e.g., 50 goroutines)
- [x] Verify all requests complete successfully with correct responses
- [x] Verify no connection leaks or race conditions occur

### 8. Integration Test: Configuration from Environment
- [x] Write a test that verifies `LISTEN_ADDR`, `BACKEND_URL`, `MAX_IDLE_CONNS`, `MAX_IDLE_CONNS_PER_HOST` environment variables are correctly parsed
- [x] Test default values when environment variables are not set

### 9. Integration Test: Large Payload
- [x] Generate a large payload (e.g., 10MB)
- [x] Send it through the proxy to a test backend
- [x] Verify the backend receives and responds correctly to the large payload
- [x] Verify the proxy doesn't hang or crash

### 10. Refactor / Create Test Helper Utilities
- [x] Create helper functions to spin up a test backend server (`httptest.NewServer`)
- [x] Create a helper to start the proxy pointing to a test backend URL
- [x] Create a helper to cleanly shut down both servers after each test
- [x] Use `t.Cleanup()` and `t.Parallel()` appropriately

### 11. Add Makefile / CI Integration
- [x] Verify `make test` runs the new integration tests (it already runs `go test ./...`)
- [x] Add a CI workflow in `.github/` to run tests on pull requests (if not already present)
- [x] Optionally add a `test-integration` target in the makefile for clarity

### 12. Documentation & Final Review
- [x] Add a `Testing` section to `README.md` explaining how to run the integration tests
- [x] Run `go vet ./...` and `golangci-lint run ./...` to ensure no lint issues
- [x] Run `go mod tidy` to clean up dependencies
- [x] Verify all tests pass with `go test -v ./...`
- [x] Review the PR diff before submitting
