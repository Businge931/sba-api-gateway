APP = sba-api-gateway
GOBASE = $(shell pwd)
GOBIN = $(GOBASE)/build/bin
LINT_PATH = $(GOBASE)/build/lint
MAIN_APP = $(GOBASE)/cmd

.PHONY: help deps build run lint lint-fix install-golangci test test-unit test-integration test-e2e docker-build docker-run docker-stop ci-pipeline

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Fetch required dependencies
	go mod tidy -compat=1.22
	go mod download

build: ## Build the application
	go build -o $(GOBIN)/$(APP) $(MAIN_APP)

run: build ## Build and run program
	cd $(MAIN_APP) && go run .

lint: install-golangci ## Linter for developers
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml

lint-fix: ## Fix linting issues automatically where possible
	$(LINT_PATH)/golangci-lint run --timeout=5m -c .golangci.yml --fix

install-golangci: ## Install the correct version of lint
	@GOBIN=$(LINT_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1

# Testing targets
test: test-unit test-integration ## Run all tests except E2E tests

test-unit: ## Run unit tests
	go test -v ./tests/unit/...

test-integration: ## Run integration tests with mocks (no containers)
	go test -v ./tests/integration/api_gateway_test.go

test-e2e: ## Run end-to-end tests with testcontainers
	@echo "Running E2E tests with testcontainers (requires Docker daemon)"
	@echo "Checking if Docker daemon is running..."
	docker info > /dev/null 2>&1 || (echo "Docker daemon is not running. Please start it first." && exit 1)
	go test -v ./tests/e2e/testcontainers_test.go

# Docker commands
docker-build: ## Build the Docker image
	docker build -t $(APP) .

docker-run: docker-build ## Run the Docker container
	docker run -d -p 8080:8080 --name $(APP)-container $(APP)

docker-stop: ## Stop and remove the Docker container
	docker stop $(APP)-container || true
	docker rm $(APP)-container || true

# CI/CD pipeline simulation
ci-pipeline: deps lint test-unit test-integration docker-build test-e2e ## Run the full CI pipeline locally

