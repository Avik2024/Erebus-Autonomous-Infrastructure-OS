.PHONY: build run test docker

build:
	go build -o bin/erebus ./cmd/erebusd

run:
	go run ./cmd/erebusd

test:
	go test ./...

docker:
	docker build -t avik2024/erebus-backend:latest .