#!/bin/bash

# Build script for shortener project
# Supports multi-architecture and cross-platform builds

set -e

PROJECT_NAME="shortener-server"
CLI_NAME="shortener"
DIST_DIR="dist"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build configuration
export CGO_ENABLED=0
export GOFLAGS="-trimpath"

# Version information
VERSION=$(cat version.txt 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Supported platforms (matching goreleaser config)
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/loong64"
    "linux/riscv64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/loong64"
    "windows/riscv64"
)

print_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  build              Build for local platform"
    echo "  build-all          Build for all platforms"
    echo "  build-server       Build server for all platforms"
    echo "  build-cli          Build CLI for all platforms"
    echo "  build-platform     Build for specific platform"
    echo "  package            Package all builds"
    echo "  clean              Clean build artifacts"
    echo "  platforms          Show supported platforms"
    echo "  version            Show version information"
    echo "  help               Show this help"
    echo ""
    echo "Options for build-platform:"
    echo "  GOOS=<os> GOARCH=<arch> $0 build-platform"
    echo "  Example: GOOS=linux GOARCH=amd64 $0 build-platform"
}

print_version() {
    echo -e "${BLUE}Version Information:${NC}"
    echo "  Version: ${VERSION}"
    echo "  Commit:  ${COMMIT}"
    echo "  Date:    ${DATE}"
}

print_platforms() {
    echo -e "${BLUE}Supported Platforms:${NC}"
    for platform in "${PLATFORMS[@]}"; do
        echo "  ${platform}"
    done
}

prepare_build() {
    echo -e "${YELLOW}Preparing build environment...${NC}"
    go mod tidy
    go generate ./... 2>/dev/null || echo "No generate tasks found, continuing..."
    mkdir -p "${DIST_DIR}"
}

build_local() {
    echo -e "${BLUE}Building for local platform...${NC}"
    prepare_build
    
    # Build server
    echo "Building server..."
    go build -ldflags "${LDFLAGS}" -o "${PROJECT_NAME}" ./main.go
    echo -e "${GREEN}✓ Built ${PROJECT_NAME}${NC}"
    
    # Build CLI
    echo "Building CLI..."
    go build -ldflags "${LDFLAGS}" -o "${CLI_NAME}" ./cmd/shortener/main.go
    echo -e "${GREEN}✓ Built ${CLI_NAME}${NC}"
}

build_server_all() {
    echo -e "${BLUE}Building server for all platforms...${NC}"
    prepare_build
    
    local success=0
    local total=0
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r goos goarch <<< "${platform}"
        output="${DIST_DIR}/${PROJECT_NAME}-server-${goos}-${goarch}"
        
        if [ "${goos}" = "windows" ]; then
            output="${output}.exe"
        fi
        
        echo "Building server for ${goos}/${goarch}..."
        total=$((total + 1))
        
        if GOOS="${goos}" GOARCH="${goarch}" go build -ldflags "${LDFLAGS}" -o "${output}" ./main.go; then
            echo -e "${GREEN}✓ Built ${output}${NC}"
            success=$((success + 1))
        else
            echo -e "${RED}✗ Failed to build server for ${goos}/${goarch}${NC}"
        fi
    done
    
    echo -e "${BLUE}Server build summary: ${success}/${total} successful${NC}"
}

build_cli_all() {
    echo -e "${BLUE}Building CLI for all platforms...${NC}"
    prepare_build
    
    local success=0
    local total=0
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r goos goarch <<< "${platform}"
        output="${DIST_DIR}/${CLI_NAME}-${goos}-${goarch}"
        
        if [ "${goos}" = "windows" ]; then
            output="${output}.exe"
        fi
        
        echo "Building CLI for ${goos}/${goarch}..."
        total=$((total + 1))
        
        if GOOS="${goos}" GOARCH="${goarch}" go build -ldflags "${LDFLAGS}" -o "${output}" ./cmd/shortener/main.go; then
            echo -e "${GREEN}✓ Built ${output}${NC}"
            success=$((success + 1))
        else
            echo -e "${RED}✗ Failed to build CLI for ${goos}/${goarch}${NC}"
        fi
    done
    
    echo -e "${BLUE}CLI build summary: ${success}/${total} successful${NC}"
}

build_all() {
    echo -e "${BLUE}Building all components for all platforms...${NC}"
    build_server_all
    build_cli_all
}

build_platform() {
    if [ -z "${GOOS}" ] || [ -z "${GOARCH}" ]; then
        echo -e "${RED}Error: GOOS and GOARCH must be set${NC}"
        echo "Usage: GOOS=<os> GOARCH=<arch> $0 build-platform"
        echo "Example: GOOS=linux GOARCH=amd64 $0 build-platform"
        exit 1
    fi
    
    echo -e "${BLUE}Building for ${GOOS}/${GOARCH}...${NC}"
    prepare_build
    
    # Build server
    server_output="${DIST_DIR}/${PROJECT_NAME}-server-${GOOS}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        server_output="${server_output}.exe"
    fi
    
    if GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags "${LDFLAGS}" -o "${server_output}" ./main.go; then
        echo -e "${GREEN}✓ Built server: ${server_output}${NC}"
    else
        echo -e "${RED}✗ Failed to build server for ${GOOS}/${GOARCH}${NC}"
    fi
    
    # Build CLI
    cli_output="${DIST_DIR}/${CLI_NAME}-${GOOS}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        cli_output="${cli_output}.exe"
    fi
    
    if GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags "${LDFLAGS}" -o "${cli_output}" ./cmd/shortener/main.go; then
        echo -e "${GREEN}✓ Built CLI: ${cli_output}${NC}"
    else
        echo -e "${RED}✗ Failed to build CLI for ${GOOS}/${GOARCH}${NC}"
    fi
}

package_all() {
    echo -e "${BLUE}Packaging all builds...${NC}"
    
    if [ ! -d "${DIST_DIR}" ]; then
        echo -e "${RED}Error: No builds found. Run build-all first.${NC}"
        exit 1
    fi
    
    mkdir -p "${DIST_DIR}/packages"
    
    cd "${DIST_DIR}"
    for file in ${PROJECT_NAME}-server-* ${CLI_NAME}-*; do
        if [ -f "${file}" ]; then
            # Extract platform info
            if [[ "${file}" == *"-server-"* ]]; then
                platform=$(echo "${file}" | sed "s/${PROJECT_NAME}-server-//" | sed 's/\.exe$//')
                name="${PROJECT_NAME}-server-${platform}"
            else
                platform=$(echo "${file}" | sed "s/${CLI_NAME}-//" | sed 's/\.exe$//')
                name="${CLI_NAME}-${platform}"
            fi
            
            # Create temp directory
            mkdir -p "temp/${name}"
            cp "${file}" "temp/${name}/"
            
            # Copy related files
            if [[ "${file}" == *"-server-"* ]]; then
                cp ../LICENSE "temp/${name}/" 2>/dev/null || true
                cp ../README.md "temp/${name}/" 2>/dev/null || true
                cp ../config/config.toml "temp/${name}/" 2>/dev/null || true
            else
                cp ../LICENSE "temp/${name}/" 2>/dev/null || true
                cp ../cmd/shortener/README.md "temp/${name}/" 2>/dev/null || true
            fi
            
            # Package
            if [[ "${platform}" == *"windows"* ]]; then
                (cd temp && zip -r "../packages/${name}.zip" "${name}/")
            else
                tar -czf "packages/${name}.tar.gz" -C temp "${name}"
            fi
            
            echo -e "${GREEN}✓ Packaged ${name}${NC}"
        fi
    done
    
    rm -rf temp
    cd ..
    echo -e "${BLUE}All packages created in ${DIST_DIR}/packages/${NC}"
}

clean_build() {
    echo -e "${YELLOW}Cleaning build artifacts...${NC}"
    rm -f "${PROJECT_NAME}"
    rm -f "${CLI_NAME}"
    rm -rf "${DIST_DIR}"
    echo -e "${GREEN}✓ Cleaned up successfully${NC}"
}

# Main script logic
case "${1:-help}" in
    "build")
        build_local
        ;;
    "build-all")
        build_all
        ;;
    "build-server")
        build_server_all
        ;;
    "build-cli")
        build_cli_all
        ;;
    "build-platform")
        build_platform
        ;;
    "package")
        package_all
        ;;
    "clean")
        clean_build
        ;;
    "platforms")
        print_platforms
        ;;
    "version")
        print_version
        ;;
    "help"|*)
        print_usage
        ;;
esac