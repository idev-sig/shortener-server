# 更新日志

## [0.3.0] - 2025-10-23

### 重大变更 (Breaking Changes)
- **模块名称变更**: 项目模块名从 `go.xoder.cn/shortener` 更改为 `go.xoder.cn/shortener-server`
  - 影响所有 Go 源文件的 import 路径
  - 需要更新 `go.mod` 和所有依赖此模块的项目

#### API 字段重命名
- **短链接相关字段**
  - `code` → `short_code` (短链接代码)
  - `describe` → `description` (描述信息)
  
- **分页相关字段**
  - `page_size` → `per_page` (每页数量)
  - `current_count` → `count` (当前页条目数)
  - `total_items` → `total` (总条目数)
  
- **历史记录字段**
  - `accessed_time` → `accessed_at` (访问时间)
  - `created_time` → `created_at` (创建时间)

- **错误响应字段**
  - `errcode` → `error_code` (错误代码)
  - `errinfo` → `error_message` (错误信息)

### 新增功能
- 新增 `/ping` 健康检查接口
- 新增 `docker/` 目录，整合 Docker 相关配置
- 新增 `config.toml` 到项目根目录

### 改进
- 更新 OpenAPI 规范版本至 3.1.2
- 优化 API 响应字段命名，更符合 RESTful 规范
- 统一时间字段命名格式（使用 `_at` 后缀）
- 改进错误信息输出格式
- 优化分页响应结构，更清晰易懂
- 更新 CLI 工具以支持新的 API 字段
- 改进 README 文档

### 修复
- 修复所有时间格式为 ISO 8601 格式（`1970-01-01T00:00:00Z`）
- 修复数据库字段映射问题

### 删除
- 删除 `config/config.toml`（已移至根目录）
- 删除 `deploy/docker/` 目录（已迁移至 `docker/`）

### 文件变更统计
- 31 个文件修改
- 新增 714 行
- 删除 912 行

### 数据库迁移

#### 修复时间格式 (urls 表)
```sqlite
-- 1. 修复所有时间（支持任意合法前缀）
UPDATE urls
SET
    created_at = strftime('%Y-%m-%dT%H:%M:%SZ', datetime(substr(created_at, 1, 19)), 'utc'),
    updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', datetime(substr(updated_at, 1, 19)), 'utc');

-- 2. 把 created_at 为 1970 的，补成 updated_at
UPDATE urls
SET created_at = updated_at
WHERE created_at = '1970-01-01T00:00:00Z';

-- 3. 验证
SELECT 
    id,
    created_at AS "创建时间",
    updated_at AS "更新时间"
FROM urls 
ORDER BY id;
```

#### 修复时间格式 (histories 表)
```sqlite
-- 1. 修复所有时间（支持任意合法前缀）
UPDATE histories
SET
    created_at = strftime('%Y-%m-%dT%H:%M:%SZ', datetime(substr(created_at, 1, 19)), 'utc'),
    accessed_at = strftime('%Y-%m-%dT%H:%M:%SZ', datetime(substr(accessed_at, 1, 19)), 'utc');

-- 2. 把 created_at 为 1970 的，补成 updated_at
UPDATE histories
SET created_at = accessed_at
WHERE created_at = '1970-01-01T00:00:00Z';

-- 3. 验证
SELECT 
    id,
    created_at AS "创建时间",
    accessed_at AS "访问时间"
FROM histories 
ORDER BY id;
```

#### 重命名字段 (urls 表)
```sqlite
ALTER TABLE urls RENAME COLUMN describe TO description;
```

### 升级注意事项
⚠️ 此版本包含破坏性变更，升级前请注意：
1. API 请求和响应字段名称已更改，需要更新客户端代码
2. 建议先在测试环境验证兼容性
3. 执行数据库迁移脚本以修复时间格式和字段名称