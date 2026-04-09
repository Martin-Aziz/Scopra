SHELL := /bin/bash

.PHONY: up down logs backend-test cli-test frontend-test test lint fmt

up:
	docker compose up --build

down:
	docker compose down -v

logs:
	docker compose logs -f

backend-test:
	cd backend && go test ./...

cli-test:
	cd cli && cargo test

frontend-test:
	cd frontend && npm test

test: backend-test cli-test frontend-test

lint:
	cd backend && go vet ./...
	cd cli && cargo clippy -- -D warnings
	cd frontend && npm run lint

fmt:
	cd backend && gofmt -w ./cmd ./src ./tests
	cd cli && cargo fmt
	cd frontend && npm run format
