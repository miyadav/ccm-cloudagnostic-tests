# Copyright 2024 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build configuration
BINARY_NAME=test-runner
E2E_BINARY_NAME=e2e-test-runner
BUILD_DIR=bin
VERSION?=$(shell git describe --tags --always --dirty)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Test configuration
TEST_TIMEOUT=10m
TEST_VERBOSE=-v
TEST_COVERAGE=coverage.out

# Go configuration
GO=go
GOLANGCI_LINT=golangci-lint
GOLANGCI_LINT_VERSION=v1.55.2

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the test runner binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/test-runner/main.go
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-e2e
build-e2e: ## Build the e2e test runner binary
	@echo "Building $(E2E_BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(E2E_BINARY_NAME) cmd/e2e-test-runner/main.go
	@echo "Binary built: $(BUILD_DIR)/$(E2E_BINARY_NAME)"

.PHONY: build-all
build-all: ## Build binaries for all supported platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/test-runner/main.go
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 cmd/test-runner/main.go
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 cmd/test-runner/main.go
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 cmd/test-runner/main.go
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe cmd/test-runner/main.go
	@echo "Binaries built in $(BUILD_DIR)/"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(TEST_COVERAGE)
	@go clean -cache
	@echo "Clean complete"

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./pkg/testing/...

.PHONY: test-e2e
test-e2e: build-e2e ## Run e2e tests with mock provider
	@echo "Running e2e tests with mock provider..."
	./$(BUILD_DIR)/$(E2E_BINARY_NAME) --provider mock --suite all --verbose

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./pkg/testing/... -run "^Test.*Unit"

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./pkg/testing/integration_test.go

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -coverprofile=$(TEST_COVERAGE) ./pkg/testing/...
	$(GO) tool cover -html=$(TEST_COVERAGE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GO) test -bench=. -benchmem ./pkg/testing/...

.PHONY: test-runner
test-runner: build ## Run the test runner with default settings
	@echo "Running test runner..."
	./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: test-runner-loadbalancer
test-runner-loadbalancer: build ## Run load balancer test suite
	@echo "Running load balancer test suite..."
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=loadbalancer -verbose

.PHONY: test-runner-nodes
test-runner-nodes: build ## Run node management test suite
	@echo "Running node management test suite..."
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=nodes -verbose

.PHONY: test-runner-routes
test-runner-routes: build ## Run route management test suite
	@echo "Running route management test suite..."
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=routes -verbose

.PHONY: test-runner-all
test-runner-all: build ## Run all test suites
	@echo "Running all test suites..."
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=all -verbose

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi
	$(GOLANGCI_LINT) run --fix

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

.PHONY: vendor
vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	$(GO) mod vendor

.PHONY: install
install: build ## Install the binary to GOPATH
	@echo "Installing binary..."
	$(GO) install $(LDFLAGS) ./cmd/test-runner

.PHONY: uninstall
uninstall: ## Remove the binary from GOPATH
	@echo "Uninstalling binary..."
	@rm -f $$(go env GOPATH)/bin/$(BINARY_NAME)

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ccm-cloudagnostic-tests:$(VERSION) .
	docker tag ccm-cloudagnostic-tests:$(VERSION) ccm-cloudagnostic-tests:latest

.PHONY: docker-run
docker-run: ## Run tests in Docker container
	@echo "Running tests in Docker container..."
	docker run --rm ccm-cloudagnostic-tests:latest

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi ccm-cloudagnostic-tests:$(VERSION) ccm-cloudagnostic-tests:latest 2>/dev/null || true

.PHONY: release
release: clean build-all ## Create release artifacts
	@echo "Creating release artifacts..."
	@mkdir -p release
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf ../release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip ../release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release artifacts created in release/"

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "All checks passed!"

.PHONY: ci
ci: deps check test-coverage ## Run CI pipeline
	@echo "CI pipeline completed!"

.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GO) mod download
	@if ! command -v $(GOLANGCI_LINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi
	@echo "Development environment setup complete!"

.PHONY: examples
examples: build ## Run example tests
	@echo "Running example tests..."
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=loadbalancer -verbose
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=nodes -verbose
	./$(BUILD_DIR)/$(BINARY_NAME) -suite=routes -verbose

.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc > /dev/null; then \
		echo "Starting godoc server on http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: update-deps
update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Go version: $(shell go version)"
	@echo "OS/Arch: $(GOOS)/$(GOARCH)"
