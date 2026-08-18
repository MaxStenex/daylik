.PHONY: deps run down migrate lint fmt test gen mocks

REDOCLY_VERSION := 1.34.5

gen: mocks
	npx --yes @redocly/cli@$(REDOCLY_VERSION) bundle api/openapi.yaml -o .build/openapi.bundled.yaml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types -package api -o internal/generated/api/api.gen.go .build/openapi.bundled.yaml

mocks:
	go run github.com/vektra/mockery/v3
	go mod tidy

deps:
	docker-compose up -d postgres

run:
	go run ./cmd/service

migrate:
	go run ./cmd/migrations

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

test:
	go test -race ./...

down:
	docker-compose down
