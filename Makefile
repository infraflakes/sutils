VERSION ?= $(shell git describe --tags --always --dirty --first-parent 2>/dev/null || echo "dev")

.PHONY: all build run fmt lint clean

all: build

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/sn .

run:
	go run . --

fmt:
	go fmt ./...

lint:
	go vet ./...
	go fix ./...
	golangci-lint run ./...

clean:
	go clean
	rm -f sn
	rm -rf result
