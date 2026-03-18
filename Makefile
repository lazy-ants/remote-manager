VERSION ?= dev
BINARY := remote-manager
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: help build test lint docker clean release

help:
	@echo "Available targets:"
	@echo "  build   - Build the binary"
	@echo "  test    - Run all tests"
	@echo "  lint    - Run golangci-lint"
	@echo "  docker  - Build Docker image"
	@echo "  clean   - Remove build artifacts"
	@echo "  release - Build with goreleaser"

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) .

test:
	go test ./...

lint:
	golangci-lint run

docker:
	docker build -t $(BINARY) .

clean:
	rm -f $(BINARY)

release:
	goreleaser release --snapshot --clean
