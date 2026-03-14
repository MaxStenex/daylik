.PHONY: deps run down migrate lint fmt test

deps:
	docker-compose up -d postgres

run:
	go run ./cmd/service

migrate:
	go run ./cmd/migrations

lint:
	golangci-lint run ./...

fmt:
	gofmt -w -l $(shell find . -name "*.go" -not -path "./vendor/*")

test:
	go test -race ./...

down:
	docker-compose down
