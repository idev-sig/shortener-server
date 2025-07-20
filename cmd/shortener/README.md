# shorten

短网址 CLI 管理工具

## 使用帮助
### 安装
```bash
go install go.bdev.cn/shortener/cmd/shortener@latest
```

### 初始化
```bash
shortener init

# 初始化时同时设置 KEY 和 URL
shortener init --key KEY --url "http://127.0.0.1:8080"
```
创建配置文件 `$HOME/.config/shortener/config.toml`
```toml
# 密钥
key = ''
# 短网址服务器
url = 'http://127.0.0.1:8080'
```
亦可通过全局变量设置
```bash
SHORTENER_KEY=""
SHORTENER_URL="http://127.0.0.1:8080"
```
临时使用可通过传参方式
```bash
shortener --key KEY --url "http://127.0.0.1:8080"
```

### 更多帮助
```bash
shortener help

# 子命令帮助
# shortener help list
```

```bash
Short URL management CLI tool

Usage:
  shortener [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a short link
  delete      Delete a short code
  env         Print environment variables
  get         Get a short link
  help        Help about any command
  init        Initialize configuration
  list        List all short links
  update      Update a short code

Flags:
  -h, --help         help for shortener
  -k, --key string   API KEY
  -u, --url string   API URL

Use "shortener [command] --help" for more information about a command.
```