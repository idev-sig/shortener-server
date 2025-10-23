# 短网址

一个超简单的短网址管理平台。

> **本项目已停止维护。**   
> **已使用 Rust 语言重构，新项目：<https://github.com/jetsung/shortener>**

- **配置前端：[shortener-frontend](https://git.jetsung.com/idev/shortener-frontend)**
- **命令行工具：[shortener](./cmd/shortener/README.md)**

## 命令行
```bash
go install go.xoder.cn/shortener-server/cmd/shortener@main
```

## [Docker](./docker/README.md)

> **版本：** `latest`, `main`, <`TAG`>

| Registry                                                                                   | Image                                                  |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------ |
| [**Docker Hub**](https://hub.docker.com/r/idevsig/shortener-server/)                                | `idevsig/shortener-server`                                    |
| [**GitHub Container Registry**](https://github.com/idev-sig/shortener-server/pkgs/container/shortener-server) | `ghcr.io/idev-sig/shortener-server`                            |
| **Tencent Cloud Container Registry（SG）**                                                       | `sgccr.ccs.tencentyun.com/idevsig/shortener-server`             |
| **Aliyun Container Registry（GZ）**                                                              | `registry.cn-guangzhou.aliyuncs.com/idevsig/shortener-server` |

## 开发

### 1. 拉取代码
```bash
git clone https://git.jetsung.com/idev/shortener-server.git
cd shortener-server
```

### 2. 修改配置
```bash
mkdir -p config/dev
cp config/config.toml config/dev/

# 修改开发环境的配置文件
vi config/dev/config.toml
```

### 3. 运行
```bash
go run .

# 或指定配置文件路径
go run . -config /path/to/config.toml
```

### 4. 构建
```bash
go build

# 支持 GoReleaser 方式构建
goreleaser release --snapshot --clean
```

### 5. 使用自定义配置文件
```bash
# 使用 -config 参数指定配置文件路径
./shortener-server -config /path/to/config.toml

# 或使用默认配置文件搜索路径（按优先级）：
# 1. ./config/dev/config.toml
# 2. ./config/prod/config.toml
# 3. ./config/config.toml
# 4. ./config.toml
./shortener-server
```

### 更多功能
```bash
just --list
```

## 文档

### Linux 部署

- 配置 Nginx 反向代理（**若使用管理界面作为入口域名，可忽略此步**）
    <details>
    <summary>点击展开/折叠</summary>

    ```nginx
    # 对接 API
    location /api/ {
        proxy_pass   http://127.0.0.1:8080/api/;

        client_max_body_size  1024m;
        proxy_set_header Host $host:$server_port;

        proxy_set_header X-Real-Ip $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;  # 透传 HTTPS 协议标识
        proxy_set_header X-Forwarded-Ssl on;         # 明确 SSL 启用状态

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_connect_timeout 99999;
    }
    ```
    </details>

1. 下载发行版的安装包：[`deb` / `rpm`](https://github.com/idev-sig/shortener-server/releases)
2. 安装
    ```bash
    # deb 安装包
    dpkg -i shortener-server_<VERSION>_linux_amd64.deb

    # rpm 安装包
    rpm -ivh shortener-server_<VERSION>_linux_amd64.rpm
    ```
3. 配置文件 `config.toml`
4. 启动
    ```bash
    systemctl start shortener-server
    systemctl enable shortener-server
    ```
5. 使用自定义配置文件（可选）
    ```bash
    # 直接运行时指定配置文件
    /usr/bin/shortener-server -config /path/to/config.toml
    
    # 或修改 systemd 服务文件
    # 编辑 /etc/systemd/system/shortener-server.service
    # 在 ExecStart 行添加 -config 参数
    # ExecStart=/usr/bin/shortener-server -config /etc/shortener-server/config.toml
    ```

**若需要前端管理平台，需要使用 [shortener-frontend](https://github.com/idev-sig/shortener-frontend/releases)** 。
1. 下载并解压到指定目录

### Docker 部署
1. 配置文件 `config.toml`
2. 若需要使用缓存，需要配置 `valkey` 缓存
    1. 取消 `compose.yml` 中的 `valkey` 配置的注释。
    2. 修改配置文件 `config.toml` 中的 `cache.enabled` 字段为 `true`。
    3. 修改配置文件 `config.toml` 中的 `cache.type` 字段为 `valkey`。
3. 若需要 IP 数据，需要配置 `ip2region` 数据库
    1. 下载 [ip2region](https://github.com/lionsoul2014/ip2region/blob/master/data/) [**V4**](https://github.com/lionsoul2014/ip2region/raw/refs/heads/master/data/ip2region_v4.xdb)、[**V6**](https://github.com/lionsoul2014/ip2region/raw/refs/heads/master/data/ip2region_v6.xdb) ，保存至 `./data/ip2region.xdb`。
    2. 修改配置文件 `config.toml` 中的 `geoip.enabled` 字段为 `true`，`path` 修改为 `xdb` 的相对路径，`version` 字段的值为 `"4"`或者 `"6"`。
4. 启动
    ```bash
    docker compose up -d
    ```
5. 配置 Nginx 反向代理
    <details>
    <summary>点击展开/折叠</summary>

    ```nginx
    # 前端配置
    location / {
        proxy_pass   http://127.0.0.1:8080;

        client_max_body_size  1024m;
        proxy_set_header Host $host:$server_port;

        proxy_set_header X-Real-Ip $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;  # 透传 HTTPS 协议标识
        proxy_set_header X-Forwarded-Ssl on;         # 明确 SSL 启用状态

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_connect_timeout 99999;
    }
    ```
    </details>

## 仓库镜像

- [MyCode](https://git.jetsung.com/idev/shortener-server)
- [Framagit](https://framagit.org/idev/shortener-server)
- [GitCode](https://gitcode.com/idev/shortener-server)
- [GitHub](https://github.com/idev-sig/shortener-server)
