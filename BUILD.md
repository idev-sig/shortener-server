# 构建指南

本项目提供了多种构建方式，支持多架构和跨平台构建。

## 构建工具

### 1. Just (推荐)

Just 是一个现代化的命令运行器，语法简洁易用。

```bash
# 安装 just (如果尚未安装)
# macOS
brew install just
# Ubuntu/Debian
sudo apt install just
# 或使用 cargo
cargo install just

# 查看所有可用命令
just help

# 本地构建
just build              # 构建服务器
just build-cli          # 构建CLI
just                    # 构建服务器和CLI (默认)

# 多平台构建
just build-all          # 构建所有平台的服务器和CLI
just build-all-server   # 构建所有平台的服务器
just build-all-cli      # 构建所有平台的CLI

# 特定平台构建
just build-platform linux amd64
just build-platform windows amd64
just build-platform darwin arm64

# 打包
just package-all        # 打包所有构建产物

# 清理
just clean              # 清理构建产物

# 查看支持的平台
just show-platforms

# 查看版本信息
just show-version
```

### 2. Make

传统的 Make 工具，兼容性好。

```bash
# 本地构建
make build              # 构建服务器
make build-cli          # 构建CLI
make                    # 构建服务器和CLI (默认)

# 多平台构建
make build-all          # 构建所有平台的服务器和CLI
make build-all-server   # 构建所有平台的服务器
make build-all-cli      # 构建所有平台的CLI

# 特定平台构建
make build-platform GOOS=linux GOARCH=amd64
make build-platform GOOS=windows GOARCH=amd64
make build-platform GOOS=darwin GOARCH=arm64

# 打包
make package-all        # 打包所有构建产物

# 清理
make clean              # 清理构建产物

# 查看支持的平台
make show-platforms

# 查看版本信息
make show-version

# 查看帮助
make help
```

### 3. Shell 脚本

直接使用 shell 脚本构建。

```bash
# 给脚本执行权限
chmod +x build.sh

# 本地构建
./build.sh build

# 多平台构建
./build.sh build-all
./build.sh build-server
./build.sh build-cli

# 特定平台构建
GOOS=linux GOARCH=amd64 ./build.sh build-platform
GOOS=windows GOARCH=amd64 ./build.sh build-platform
GOOS=darwin GOARCH=arm64 ./build.sh build-platform

# 打包
./build.sh package

# 清理
./build.sh clean

# 查看支持的平台
./build.sh platforms

# 查看版本信息
./build.sh version

# 查看帮助
./build.sh help
```

## 支持的平台

### 操作系统
- Linux
- macOS (Darwin)
- Windows

### 架构
- amd64 (x86_64)
- arm64 (aarch64)
- loong64 (LoongArch)
- riscv64 (RISC-V 64位)

### 支持的平台组合
- linux/amd64
- linux/arm64
- linux/loong64
- linux/riscv64
- darwin/amd64
- darwin/arm64
- windows/amd64
- windows/loong64
- windows/riscv64

## 构建产物

### 目录结构

```
dist/
├── shortener-server-linux-amd64
├── shortener-server-linux-arm64
├── shortener-server-linux-loong64
├── shortener-server-linux-riscv64
├── shortener-server-darwin-amd64
├── shortener-server-darwin-arm64
├── shortener-server-windows-amd64.exe
├── shortener-server-windows-loong64.exe
├── shortener-server-windows-riscv64.exe
├── shortener-linux-amd64
├── shortener-linux-arm64
├── shortener-linux-loong64
├── shortener-linux-riscv64
├── shortener-darwin-amd64
├── shortener-darwin-arm64
├── shortener-windows-amd64.exe
├── shortener-windows-loong64.exe
├── shortener-windows-riscv64.exe
└── packages/
    ├── shortener-server-linux-amd64.tar.gz
    ├── shortener-server-darwin-arm64.tar.gz
    ├── shortener-server-windows-amd64.zip
    ├── shortener-linux-amd64.tar.gz
    ├── shortener-darwin-arm64.tar.gz
    └── shortener-windows-amd64.zip
```

### 打包内容

**服务器包含:**
- 可执行文件
- LICENSE
- README.md
- config.toml (如果存在)

**CLI包含:**
- 可执行文件
- LICENSE
- CLI README.md (如果存在)

## 版本信息

构建时会自动注入版本信息：
- `main.version`: 从 `version.txt` 读取，或默认为 "dev"
- `main.commit`: Git commit hash (短格式)
- `main.date`: 构建时间 (UTC)

## 环境变量

构建过程中使用的环境变量：
- `CGO_ENABLED=0`: 禁用 CGO，确保静态链接
- `GOFLAGS=-trimpath`: 移除构建路径信息

## 构建标志

使用的 Go 构建标志：
- `-s`: 去除符号表
- `-w`: 去除调试信息
- `-X main.version=...`: 注入版本信息
- `-X main.commit=...`: 注入提交信息
- `-X main.date=...`: 注入构建时间

## 示例用法

### 快速开始

```bash
# 使用 just (推荐)
just build-all package-all

# 使用 make
make build-all package-all

# 使用脚本
./build.sh build-all
./build.sh package
```

### 构建特定平台

```bash
# Linux ARM64
just build-platform linux arm64

# Windows AMD64
just build-platform windows amd64

# macOS ARM64 (Apple Silicon)
just build-platform darwin arm64
```

### 仅构建服务器或CLI

```bash
# 仅构建所有平台的服务器
just build-all-server

# 仅构建所有平台的CLI
just build-all-cli
```

## 故障排除

### 常见问题

1. **构建失败**: 确保 Go 版本 >= 1.21
2. **权限错误**: 确保脚本有执行权限 `chmod +x build.sh`
3. **平台不支持**: 检查 Go 是否支持目标平台 `go tool dist list`

### 调试构建

```bash
# 查看详细构建信息
GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w" -o dist/test ./main.go

# 检查二进制文件信息
file dist/shortener-server-linux-amd64
ldd dist/shortener-server-linux-amd64  # Linux
otool -L dist/shortener-server-darwin-amd64  # macOS
```

## 性能优化

### 并行构建

构建脚本支持并行构建以提高速度：

```bash
# 设置并行构建数量
export GOMAXPROCS=8

# 或使用 xargs 并行构建
echo "linux/amd64 linux/arm64 darwin/amd64" | tr ' ' '\n' | xargs -P 3 -I {} bash -c 'GOOS=$(echo {} | cut -d/ -f1) GOARCH=$(echo {} | cut -d/ -f2) ./build.sh build-platform'
```

### 缓存优化

```bash
# 预热模块缓存
go mod download

# 使用构建缓存
export GOCACHE=/tmp/go-build-cache
```