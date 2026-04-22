.PHONY: deps run down migrate lint fmt test gen

REDOCLY_VERSION := 1.34.5

gen:
	npx --yes @redocly/cli@$(REDOCLY_VERSION) bundle api/openapi.yaml -o .build/openapi.bundled.yaml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types -package api -o internal/generated/api/api.gen.go .build/openapi.bundled.yaml

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
