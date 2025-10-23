# Justfile for shortener project - Multi-architecture builds

# Project configuration
project_name := "shortener-server"
cli_name := "shortener"
dist_dir := "dist"

# Build configuration
export CGO_ENABLED := "0"
export GOFLAGS := "-trimpath"

# Version information
version := `cat version.txt 2>/dev/null || echo "dev"`
commit := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
date := `date -u +"%Y-%m-%dT%H:%M:%SZ"`

# Build flags
ldflags := "-s -w -X main.version=" + version + " -X main.commit=" + commit + " -X main.date=" + date

# Supported platforms (matching goreleaser config)
platforms := "linux/amd64 linux/arm64 linux/loong64 linux/riscv64 darwin/amd64 darwin/arm64 windows/amd64 windows/loong64 windows/riscv64"

# Default target
default: build build-cli

# Build server for local platform
[group('build')]
build:
    @echo "Building {{project_name}} for local platform..."
    @go mod tidy
    @go generate ./... || echo "No generate tasks found, continuing..."
    @go build -ldflags "{{ldflags}}" -o "{{project_name}}" ./main.go
    @echo "Built {{project_name}} successfully"

# Build CLI for local platform
[group('build')]
build-cli:
    @echo "Building {{cli_name}} CLI for local platform..."
    @go mod tidy
    @go generate ./... || echo "No generate tasks found, continuing..."
    @go build -ldflags "{{ldflags}}" -o "{{cli_name}}" ./cmd/shortener/main.go
    @echo "Built {{cli_name}} successfully"

# Build server for all platforms
[unix]
[group('build')]
build-all-server:
    #!/usr/bin/env -S bash -x
    echo "Building {{project_name}} server for all platforms..."
    mkdir -p {{dist_dir}}
    go mod tidy
    go generate ./... || echo "No generate tasks found, continuing..."
    for platform in "{{platforms}}"; do
        IFS='/' read -r goos goarch <<< "$platform"
        echo "platform: $platform"
        output="{{dist_dir}}/{{project_name}}-server-${goos}-${goarch}"
        if [ "$goos" = "windows" ]; then
            output="${output}.exe"
        fi
        echo "Building server for $goos/$goarch..."
        GOOS=$goos GOARCH=$goarch go build -ldflags "{{ldflags}}" -o "$output" ./main.go
        if [ $? -eq 0 ]; then
            echo "✓ Built $output"
        else
            echo "✗ Failed to build server for $goos/$goarch"
        fi
    done
    echo "All server builds completed"

# Build CLI for all platforms
[group('build')]
build-all-cli:
    #!/usr/bin/env -S bash -x
    echo "Building {{cli_name}} CLI for all platforms..."
    mkdir -p {{dist_dir}}
    go mod tidy
    go generate ./... || echo "No generate tasks found, continuing..."
    for platform in {{platforms}}; do
        IFS='/' read -r goos goarch <<< "$platform"
        output="{{dist_dir}}/{{cli_name}}-${goos}-${goarch}"
        if [ "$goos" = "windows" ]; then
            output="${output}.exe"
        fi
        echo "Building CLI for $goos/$goarch..."
        GOOS=$goos GOARCH=$goarch go build -ldflags "{{ldflags}}" -o "$output" ./cmd/shortener/main.go
        if [ $? -eq 0 ]; then
            echo "✓ Built $output"
        else
            echo "✗ Failed to build CLI for $goos/$goarch"
        fi
    done
    @echo "All CLI builds completed"

# Build all platforms for both server and CLI
[group('build')]
build-all: build-all-server build-all-cli

# Build for specific platform
[group('build')]
build-platform GOOS GOARCH:
    #!/usr/bin/env -S bash -x
    echo "Building for {{GOOS}}/{{GOARCH}}..."
    mkdir -p {{dist_dir}}
    go mod tidy
    # Build server
    server_output="{{dist_dir}}/{{project_name}}-server-{{GOOS}}-{{GOARCH}}"
    if [ "{{GOOS}}" = "windows" ]; then
        server_output="${server_output}.exe"
    fi
    GOOS={{GOOS}} GOARCH={{GOARCH}} go build -ldflags "{{ldflags}}" -o "$server_output" ./main.go

    # Build CLI
    cli_output="{{dist_dir}}/{{cli_name}}-{{GOOS}}-{{GOARCH}}"
    if [ "{{GOOS}}" = "windows" ]; then
        cli_output="${cli_output}.exe"
    fi
    GOOS={{GOOS}} GOARCH={{GOARCH}} go build -ldflags "{{ldflags}}" -o "$cli_output" ./cmd/shortener/main.go

    echo "✓ Built server: $server_output"
    echo "✓ Built CLI: $cli_output"

# Package all builds
package-all: build-all
    #!/usr/bin/env -S bash -x
    echo "Packaging all builds..."
    mkdir -p {{dist_dir}}/packages
    cd {{dist_dir}}
    for file in {{project_name}}-server-* {{cli_name}}-*; do
        if [ -f "$file" ]; then
            # Extract platform info
            if [[ "$file" == *"-server-"* ]]; then
                platform=$(echo "$file" | sed 's/{{project_name}}-server-//' | sed 's/\.exe$//')
                name="{{project_name}}-server-${platform}"
            else
                platform=$(echo "$file" | sed 's/{{cli_name}}-//' | sed 's/\.exe$//')
                name="{{cli_name}}-${platform}"
            fi

            # Create temp directory
            mkdir -p "temp/$name"
            cp "$file" "temp/$name/"

            # Copy related files
            if [[ "$file" == *"-server-"* ]]; then
                cp ../LICENSE "temp/$name/" 2>/dev/null || true
                cp ../README.md "temp/$name/" 2>/dev/null || true
                cp ../config/config.toml "temp/$name/" 2>/dev/null || true
            else
                cp ../LICENSE "temp/$name/" 2>/dev/null || true
                cp ../cmd/shortener/README.md "temp/$name/" 2>/dev/null || true
            fi

            # Package
            if [[ "$platform" == *"windows"* ]]; then
                (cd temp && zip -r "../packages/${name}.zip" "$name/")
            else
                tar -czf "packages/${name}.tar.gz" -C temp "$name"
            fi

            echo "✓ Packaged $name"
        fi
    done
    rm -rf temp
    @echo "All packages created in {{dist_dir}}/packages/"

# Clean build artifacts
clean:
    @echo "Cleaning up..."
    @rm -f "{{project_name}}"
    @rm -f "{{cli_name}}"
    @rm -rf "{{dist_dir}}"
    @echo "Cleaned up successfully"

# Docker operations
[group('docker')]
docker-start:
    @docker compose --profile valkey -f docker/compose.yml up -d

[group('docker')]
docker-stop:
    @docker compose --profile valkey -f docker/compose.yml down

# Testing
[group('test')]
test: docker-start test-ci docker-stop

[group('test')]
test-ci:
    @echo "Running CI tests..."
    @go test ./... -v

# Code formatting and linting
fmt:
    @gofumpt -w ./
    @goimports -w -local go.bdev.cn/shortener ./

lint:
    @golangci-lint run

# Dependency management
go-mod-tidy:
    #!/usr/bin/env -S bash -x
    find . -type f -name 'go.mod' -exec dirname {} \; | sort | while read dir; do
        echo "go mod tidy in $dir"
        (cd "$dir" && \
            go get -u ./... && \
            go mod tidy)
    done

# Pre-commit checks
pre-commit:
    @echo "Running pre-commit checks..."
    @just go-mod-tidy
    @just fmt
    @just lint

# Show supported platforms
show-platforms:
    #!/usr/bin/env -S bash -x
    echo "Supported platforms:"
    for platform in {{platforms}}; do
        echo "  $platform"
    done

# Show version information
show-version:
    @echo "Version: {{version}}"
    @echo "Commit: {{commit}}"
    @echo "Date: {{date}}"

# Help
help:
    @echo "Available commands:"
    @echo "  build              - Build server for local platform"
    @echo "  build-cli          - Build CLI for local platform"
    @echo "  build-all          - Build both server and CLI for all platforms"
    @echo "  build-all-server   - Build server for all platforms"
    @echo "  build-all-cli      - Build CLI for all platforms"
    @echo "  build-platform     - Build for specific platform (usage: just build-platform linux amd64)"
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