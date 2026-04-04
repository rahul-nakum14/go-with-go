# =============================================================================
# Go User CRUD — Makefile
# =============================================================================

APP       := go-user-crud
CMD       := ./cmd/api
MODULE    := github.com/yourorg/go-user-crud
GOFLAGS   := -race
BUILD_DIR := bin

## ── Build ─────────────────────────────────────────────────────────────────────

.PHONY: build
build:			## Build the binary
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(APP) $(CMD)

.PHONY: run
run:			## Run the server locally (requires .env)
	@set -a && source .env && set +a && go run $(GOFLAGS) $(CMD)/...

.PHONY: watch
watch:			## Live-reload with 'air' (install: go install github.com/air-verse/air@latest)
	air

## ── Test ──────────────────────────────────────────────────────────────────────

.PHONY: test
test:			## Run unit tests
	go test $(GOFLAGS) ./... -count=1 -timeout 60s

.PHONY: test/cover
test/cover:		## Run tests and open HTML coverage report
	go test $(GOFLAGS) ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: test/integration
test/integration:	## Run integration tests (requires TEST_DATABASE_URL)
	go test $(GOFLAGS) ./... -tags=integration -count=1 -timeout 120s

## ── Quality ───────────────────────────────────────────────────────────────────

.PHONY: lint
lint:			## Run golangci-lint
	golangci-lint run ./...

.PHONY: vet
vet:			## Run go vet
	go vet ./...

.PHONY: fmt
fmt:			## Format all Go files
	gofmt -s -w .
	goimports -w .

.PHONY: tidy
tidy:			## Tidy go.mod / go.sum
	go mod tidy

.PHONY: check
check: fmt vet lint test	## Run all quality checks

## ── Database ──────────────────────────────────────────────────────────────────

.PHONY: migrate/up
migrate/up:		## Apply all pending migrations
	migrate -path ./migrations -database "$$DATABASE_URL" up

.PHONY: migrate/down
migrate/down:		## Roll back the last migration
	migrate -path ./migrations -database "$$DATABASE_URL" down 1

.PHONY: migrate/create
migrate/create:		## Create a new migration: make migrate/create NAME=add_foo
	migrate create -ext sql -dir ./migrations -seq $(NAME)

## ── Docker ────────────────────────────────────────────────────────────────────

.PHONY: docker/up
docker/up:		## Start Postgres + app via docker compose
	docker compose up -d

.PHONY: docker/down
docker/down:		## Stop and remove containers
	docker compose down

.PHONY: docker/build
docker/build:		## Build the Docker image
	docker build -t $(APP):latest .

## ── Code generation ───────────────────────────────────────────────────────────

.PHONY: generate
generate:		## Run go generate (mocks, swagger, etc.)
	go generate ./...

## ── Helpers ───────────────────────────────────────────────────────────────────

.PHONY: help
help:			## Show this help message
	@grep -E '^[a-zA-Z_/]+:.*?##' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
