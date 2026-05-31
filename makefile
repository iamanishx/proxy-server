APP_NAME = proxy-server
MAIN = ./cmd/main.go
BINARY = bin/$(APP_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

build:
	go build -o $(BINARY) $(MAIN)

run:
	go run $(MAIN)

test:
	go test ./...

test-integration:
	go test -v -timeout 120s -count=1 ./...

check: fmt vet tidy lint

clean:
	rm -rf $(BINARY)
