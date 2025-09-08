.PHONY: build run test docker

APP_NAME = erebusd
PKG = github.com/Avik2024/erebus/backend/internal/version

# Get current git commit and UTC date
commit := $(shell git rev-parse HEAD)
date := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
version := v0.0.1

build:
	go build -ldflags "-X $(PKG).Version=$(version) -X $(PKG).Commit=$(commit) -X $(PKG).Date=$(date)" -o bin/$(APP_NAME) ./cmd/erebusd

run: build
	./bin/$(APP_NAME)

test:
	go test ./...

docker:
	docker build -t avik2024/erebus-backend:latest .
