# 短网址

一个超简单的短网址管理平台。

**配置前端：[shortener-frontend](https://git.jetsung.com/idev/shortener-frontend)**   
**命令行工具：[shortener](./cmd/shortener/README.md)**   

## 命令行
```bash
go install go.dsig.cn/shortener/cmd/shortener@latest
```

## [Docker](./deploy/docker/README.md)

> **版本：** `latest`, `main`, <`TAG`>

| Registry                                                                                   | Image                                                  |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------ |
| [**Docker Hub**](https://hub.docker.com/r/idevsig/shortener-server/)                                | `idevsig/shortener-server`                                    |
| [**GitHub Container Registry**](https://github.com/idevsig/shortener-server/pkgs/container/shortener-server) | `ghcr.io/idevsig/shortener-server`                            |
| **Tencent Cloud Container Registry**                                                       | `ccr.ccs.tencentyun.com/idevsig/shortener-server`             |
| **Aliyun Container Registry**                                                              | `registry.cn-guangzhou.aliyuncs.com/idevsig/shortener-server` |

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
```

### 4. 构建
```bash
go build

# 支持 GoReleaser 方式构建
goreleaser release --snapshot --clean
```

### 更多功能
```bash
just --list
```

## 文档

### Linux 部署

1. 下载发行版的安装包：[`deb` / `rpm`](https://github.com/idevsig/shortener-server/releases)
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
5. 配置 Nginx 反向代理（**若使用管理界面作为入口域名，可忽略此步**）
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

**若需要前端管理平台，需要使用 [shortener-frontend](https://github.com/idevsig/shortener-frontend/releases)** 。
1. 下载并解压到指定目录
2. 配置 `nginx`：
    <details>
    <summary>点击展开/折叠</summary>
    
    ```nginx
    ...
    listen 80;

    server_name <DOMAIN>;

    index index.html;
    root /data/wwwroot/<DOMAIN>;

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

### Docker 部署
1. 配置文件 `config.toml`
2. 若需要使用缓存，需要配置 `valkey` 缓存
    1. 取消 `compose.yml` 中的 `valkey` 配置的注释。
    2. 修改配置文件 `config.toml` 中的 `cache.enabled` 字段为 `true`。
    3. 修改配置文件 `config.toml` 中的 `cache.type` 字段为 `valkey`。
3. 若需要 IP 数据，需要配置 `ip2region` 数据库
    1. 下载 [ip2region.xdb](https://github.com/lionsoul2014/ip2region/blob/master/data/ip2region.xdb) ，保存至 `./data/ip2region.xdb`。
    2. 修改配置文件 `config.toml` 中的 `geoip.enabled` 字段为 `true`。
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

## TODO

- [x] 实现全部功能接口
  - [x] `API` 权限校验
- [x] 支持数据库
  - [x] SQLite
  - [x] PostgreSQL
  - [x] MySQL
- [x] 支持缓存
  - [x] Redis
  - [x] Valkey
- [x] 制作 CLI 工具
  - [x] 添加 OpenAPI
- [x] 添加跳转请求日志记录
- [x] `CI/CD` 构建
  - [x] Docker 镜像构建与推送
- [x] 实现管理平台接口
- [x] 添加文档
- [ ] 添加测试

## 仓库镜像

- https://git.jetsung.com/idev/shortener-server
- https://framagit.org/idev/shortener-server
- https://gitcode.com/idev/shortener-server
- https://github.com/idevsig/shortener-server
