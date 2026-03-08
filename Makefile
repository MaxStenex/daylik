.PHONY: deps run down migrate lint test

deps:
	docker-compose up -d postgres

run:
	go run ./cmd/service

migrate:
	go run ./cmd/migrations

lint:
	golangci-lint run ./...

test:
	go test -race ./...

down:
	docker-compose down
