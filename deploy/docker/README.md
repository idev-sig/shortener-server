# Docker 部署

## Docker Compose
- [`compose.yml`](compose.yml)

### SQLite
```yaml
services:
  shortener:
    image: idevsig/shortener-server:latest
    container_name: shortener
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release
      - DATABASE_TYPE=sqlite
    ports:
      - "8080:8080"
    volumes:
      - ./config.toml:/app/config.toml
```

### PostgreSQL
```yaml
services:
  shortener:
    image: idevsig/shortener-server:latest
    container_name: shortener
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release
      - DATABASE_TYPE=postgres
    ports:
      - "8080:8080"
    volumes:
      - ./config.toml:/app/config.toml
    depends_on:
      - postgres

  postgres:
    image: postgres:latest
    container_name: postgres
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=shortener
```

### MySQL
```yaml
services:
  shortener:
    image: idevsig/shortener-server:latest
    container_name: shortener
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release
      - DATABASE_TYPE=mysql
    ports:
      - "8080:8080"
    volumes:
      - ./config.toml:/app/config.toml
    depends_on:
      - mysql

  mysql:
    image: mysql:latest
    container_name: mysql
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
      - MYSQL_USER=root
      - MYSQL_PASSWORD=root
      - MYSQL_ROOT_PASSWORD=root
      - MYSQL_DATABASE=shortener
```

## 运行测试

- 数据库
```bash
# sqlite test
DATABASE_TYPE=sqlite docker compose --profile sqlite up -d
DATABASE_TYPE=sqlite docker compose --profile sqlite down

# mysql test
DATABASE_TYPE=mysql docker compose --profile mysql up -d
DATABASE_TYPE=mysql docker compose --profile mysql down

# postgres test
DATABASE_TYPE=postgres docker compose --profile postgres up -d
DATABASE_TYPE=postgres docker compose --profile postgres down
```

# sqlite valkey test
```bash
DATABASE_TYPE=sqlite docker compose --profile sqlite up -d
DATABASE_TYPE=sqlite docker compose --profile sqlite down
```