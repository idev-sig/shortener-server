# Makefile for shortener project

# Project configuration
PROJECT_NAME := $(shell basename $(PWD))
CLI_NAME := shortener
DIST_DIR := dist

# Build configuration
export CGO_ENABLED := 0
export GOFLAGS := -trimpath

# Version information
VERSION := $(shell cat version.txt 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Supported platforms (matching goreleaser config)
PLATFORMS := linux/amd64 linux/arm64 linux/loong64 linux/riscv64 \
             darwin/amd64 darwin/arm64 \
             windows/amd64 windows/loong64 windows/riscv64

# Default target
.PHONY: all
all: build build-cli

# Build server for local platform
.PHONY: build
build:
	@echo "Building $(PROJECT_NAME) for local platform..."
	@go mod tidy
	@go generate ./... || echo "No generate tasks found, continuing..."
	@go build -ldflags "$(LDFLAGS)" -o "$(PROJECT_NAME)" ./main.go
	@echo "Built $(PROJECT_NAME) successfully"

# Build CLI for local platform
.PHONY: build-cli
build-cli:
	@echo "Building $(CLI_NAME) CLI for local platform..."
	@go mod tidy
	@go generate ./... || echo "No generate tasks found, continuing..."
	@go build -ldflags "$(LDFLAGS)" -o "$(CLI_NAME)" ./cmd/shortener/main.go
	@echo "Built $(CLI_NAME) successfully"

# Build server for all platforms
.PHONY: build-all-server
build-all-server:
	@echo "Building $(PROJECT_NAME) server for all platforms..."
	@mkdir -p $(DIST_DIR)
	@go mod tidy
	@go generate ./... || echo "No generate tasks found, continuing..."
	@for platform in $(PLATFORMS); do \
		goos=$$(echo $$platform | cut -d'/' -f1); \
		goarch=$$(echo $$platform | cut -d'/' -f2); \
		output="$(DIST_DIR)/$(PROJECT_NAME)-server-$$goos-$$goarch"; \
		if [ "$$goos" = "windows" ]; then \
			output="$$output.exe"; \
		fi; \
		echo "Building server for $$goos/$$goarch..."; \
		GOOS=$$goos GOARCH=$$goarch go build -ldflags "$(LDFLAGS)" -o "$$output" ./main.go; \
		if [ $$? -eq 0 ]; then \
			echo "✓ Built $$output"; \
		else \
			echo "✗ Failed to build server for $$goos/$$goarch"; \
		fi; \
	done
	@echo "All server builds completed"

# Build CLI for all platforms
.PHONY: build-all-cli
build-all-cli:
	@echo "Building $(CLI_NAME) CLI for all platforms..."
	@mkdir -p $(DIST_DIR)
	@go mod tidy
	@go generate ./... || echo "No generate tasks found, continuing..."
	@for platform in $(PLATFORMS); do \
		goos=$$(echo $$platform | cut -d'/' -f1); \
		goarch=$$(echo $$platform | cut -d'/' -f2); \
		output="$(DIST_DIR)/$(CLI_NAME)-$$goos-$$goarch"; \
		if [ "$$goos" = "windows" ]; then \
			output="$$output.exe"; \
		fi; \
		echo "Building CLI for $$goos/$$goarch..."; \
		GOOS=$$goos GOARCH=$$goarch go build -ldflags "$(LDFLAGS)" -o "$$output" ./cmd/shortener/main.go; \
		if [ $$? -eq 0 ]; then \
			echo "✓ Built $$output"; \
		else \
			echo "✗ Failed to build CLI for $$goos/$$goarch"; \
		fi; \
	done
	@echo "All CLI builds completed"

# Build all platforms for both server and CLI
.PHONY: build-all
build-all: build-all-server build-all-cli

# Package all builds
.PHONY: package-all
package-all: build-all
	@echo "Packaging all builds..."
	@mkdir -p $(DIST_DIR)/packages
	@cd $(DIST_DIR) && \
	for file in $(PROJECT_NAME)-server-* $(CLI_NAME)-*; do \
		if [ -f "$$file" ]; then \
			if echo "$$file" | grep -q "server"; then \
				platform=$$(echo "$$file" | sed 's/$(PROJECT_NAME)-server-//' | sed 's/\.exe$$//'); \
				name="$(PROJECT_NAME)-server-$$platform"; \
			else \
				platform=$$(echo "$$file" | sed 's/$(CLI_NAME)-//' | sed 's/\.exe$$//'); \
				name="$(CLI_NAME)-$$platform"; \
			fi; \
			mkdir -p "temp/$$name"; \
			cp "$$file" "temp/$$name/"; \
			if echo "$$file" | grep -q "server"; then \
				cp ../LICENSE "temp/$$name/" 2>/dev/null || true; \
				cp ../README.md "temp/$$name/" 2>/dev/null || true; \
				cp ../config/config.toml "temp/$$name/" 2>/dev/null || true; \
			else \
				cp ../LICENSE "temp/$$name/" 2>/dev/null || true; \
				cp ../cmd/shortener/README.md "temp/$$name/" 2>/dev/null || true; \
			fi; \
			if echo "$$platform" | grep -q "windows"; then \
				(cd temp && zip -r "../packages/$$name.zip" "$$name/"); \
			else \
				tar -czf "packages/$$name.tar.gz" -C temp "$$name"; \
			fi; \
			echo "✓ Packaged $$name"; \
		fi; \
	done && \
	rm -rf temp
	@echo "All packages created in $(DIST_DIR)/packages/"

# Build for specific platform
.PHONY: build-platform
build-platform:
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Usage: make build-platform GOOS=<os> GOARCH=<arch>"; \
		echo "Example: make build-platform GOOS=linux GOARCH=amd64"; \
		exit 1; \
	fi
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(DIST_DIR)
	@go mod tidy
	@server_output="$(DIST_DIR)/$(PROJECT_NAME)-server-$(GOOS)-$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then \
		server_output="$$server_output.exe"; \
	fi; \
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o "$$server_output" ./main.go; \
	cli_output="$(DIST_DIR)/$(CLI_NAME)-$(GOOS)-$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then \
		cli_output="$$cli_output.exe"; \
	fi; \
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o "$$cli_output" ./cmd/shortener/main.go; \
	echo "✓ Built server: $$server_output"; \
	echo "✓ Built CLI: $$cli_output"

# Goreleaser builds
.PHONY: build-snapshot-cli
build-snapshot-cli:
	@echo "Building $(PROJECT_NAME) CLI snapshot..."
	@goreleaser release --snapshot --clean --config .goreleaser-cli.yaml
	@echo "Built $(PROJECT_NAME) CLI snapshot successfully"

.PHONY: build-release-cli
build-release-cli:
	@echo "Building $(PROJECT_NAME) CLI release..."
	@goreleaser release --clean --config .goreleaser-cli.yaml
	@echo "Built $(PROJECT_NAME) CLI release successfully"

.PHONY: build-snapshot-server
build-snapshot-server:
	@echo "Building $(PROJECT_NAME) server snapshot..."
	@goreleaser release --snapshot --clean --config .goreleaser-server.yaml
	@echo "Built $(PROJECT_NAME) server snapshot successfully"

.PHONY: build-release-server
build-release-server:
	@echo "Building $(PROJECT_NAME) server release..."
	@goreleaser release --clean --config .goreleaser-server.yaml
	@echo "Built $(PROJECT_NAME) server release successfully"

# Prepare build artifacts
.PHONY: tidy
tidy: build
	@echo "Tidying $(PROJECT_NAME)..."
	@rm -rf "$(DIST_DIR)/$(PROJECT_NAME)"
	@mkdir -p "$(DIST_DIR)/$(PROJECT_NAME)/data"
	@cp "$(PROJECT_NAME)" "$(DIST_DIR)/$(PROJECT_NAME)/"
	@cp -f "config/config.toml" "$(DIST_DIR)/$(PROJECT_NAME)/config.toml" || echo "Warning: config.toml not found"
	@cp -f "LICENSE" "$(DIST_DIR)/$(PROJECT_NAME)/LICENSE" || echo "Warning: LICENSE not found"
	@cp -f "README.md" "$(DIST_DIR)/$(PROJECT_NAME)/README.md" || echo "Warning: README.md not found"
	@echo "Tidied $(PROJECT_NAME) successfully"

# Package build artifacts
.PHONY: package
package: tidy
	@echo "Packaging $(PROJECT_NAME)..."
	@tar -czf "$(DIST_DIR)/$(PROJECT_NAME).tar.gz" -C "$(DIST_DIR)/$(PROJECT_NAME)" .
	@echo "Packaged $(PROJECT_NAME) successfully"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f "$(PROJECT_NAME)"
	@rm -f "$(CLI_NAME)"
	@rm -rf "$(DIST_DIR)"
	@echo "Cleaned up successfully"

# Docker operations
.PHONY: docker-start
docker-start:
	@docker compose --profile valkey -f docker/compose.yml up -d

.PHONY: docker-stop
docker-stop:
	@docker compose --profile valkey -f docker/compose.yml down

# Testing
.PHONY: test
test: docker-start test-ci docker-stop

.PHONY: test-ci
test-ci:
	@echo "Running CI tests..."
	@go test ./... -v

# Code formatting and linting
.PHONY: fmt
fmt:
	@gofumpt -w ./
	@goimports -w -local go.bdev.cn/shortener ./

.PHONY: lint
lint:
	@golangci-lint run

# Dependency management
.PHONY: go-mod-tidy
go-mod-tidy:
	@find . -type f -name 'go.mod' -exec dirname {} \; | sort | while read dir; do \
		echo "go mod tidy in $$dir"; \
		(cd "$$dir" && go get -u ./... && go mod tidy); \
	done

# Pre-commit checks
.PHONY: pre-commit
pre-commit:
	@echo "Running pre-commit checks..."
	@$(MAKE) go-mod-tidy
	@$(MAKE) fmt
	@$(MAKE) lint

# Show supported platforms
.PHONY: show-platforms
show-platforms:
	@echo "Supported platforms:"
	@for platform in $(PLATFORMS); do \
		echo "  $$platform"; \
	done

# Show version information
.PHONY: show-version
show-version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build server for local platform"
	@echo "  build-cli          - Build CLI for local platform"
	@echo "  build-all          - Build both server and CLI for all platforms"
	@echo "  build-all-server   - Build server for all platforms"
	@echo "  build-all-cli      - Build CLI for all platforms"
	@echo "  build-platform     - Build for specific platform (requires GOOS and GOARCH)"
	@echo "  package-all        - Package all builds into archives"
	@echo "  clean              - Clean build artifacts"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  test               - Run tests with Docker"
	@echo "  docker-start       - Start Docker services"
	@echo "  docker-stop        - Stop Docker services"
	@echo "  show-platforms     - Show supported platforms"
	@echo "  show-version       - Show version information"
	@echo "  help               - Show this help message"