package bootstrap

import (
	"os"

	"github.com/spf13/viper"

	"go.bdev.cn/shortener/internal/shared"
	"go.bdev.cn/shortener/internal/types"
	"go.bdev.cn/shortener/internal/utils"
)

// initSharedConfig 初始化共享配置
func initSharedConfig() {
	// log.Println("shorten init")
	length := viper.GetInt("shortener.code_length")
	charset := viper.GetString("shortener.code_charset")

	if length <= 0 {
		length = 6
	}

	if charset == "" {
		charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	shared.GlobalShorten = &types.CfgShorten{
		Length:  length,
		Charset: charset,
	}

	initAPIKeyConfig()

	initUserConfig()
}

// initDefaultConfig 初始化默认配置
func initDefaultConfig() {
	// 服务器配置
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.trusted-platform", "")
	viper.SetDefault("server.site_url", "http://localhost:8080")
	viper.SetDefault("server.api_key", "")

	// 短链生成配置
	viper.SetDefault("shortener.code_length", 6)
	viper.SetDefault("shortener.code_charset", "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// 登录账号和密码
	viper.SetDefault("admin.username", "")
	viper.SetDefault("admin.password", "")

	// 数据库配置
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.log_level", 0)

	// SQLite配置
	viper.SetDefault("database.sqlite.path", "data/shortener.db")

	// PostgreSQL配置
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.user", "postgres")
	viper.SetDefault("database.postgres.password", "postgres")
	viper.SetDefault("database.postgres.database", "shortener")
	viper.SetDefault("database.postgres.sslmode", "disable")
	viper.SetDefault("database.postgres.timezone", "Asia/Shanghai")

	// MySQL配置
	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.user", "root")
	viper.SetDefault("database.mysql.password", "root")
	viper.SetDefault("database.mysql.database", "shortener")
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.parse_time", true)
	viper.SetDefault("database.mysql.loc", "Local")

	// 缓存配置
	viper.SetDefault("cache.enabled", false)
	viper.SetDefault("cache.type", "redis")
	viper.SetDefault("cache.expire", 3600)
	viper.SetDefault("cache.prefix", "shorten:")

	// Redis配置
	viper.SetDefault("cache.redis.host", "localhost")
	viper.SetDefault("cache.redis.port", 6379)
	viper.SetDefault("cache.redis.password", "")
	viper.SetDefault("cache.redis.db", 0)

	// Valkey配置
	viper.SetDefault("cache.valkey.host", "localhost")
	viper.SetDefault("cache.valkey.port", 6379)
	viper.SetDefault("cache.valkey.username", "")
	viper.SetDefault("cache.valkey.password", "")
	viper.SetDefault("cache.valkey.db", 0)

	// GeoIP配置
	viper.SetDefault("geoip.enabled", false)
	viper.SetDefault("geoip.type", "ip2region")
	viper.SetDefault("geoip.ip2region.path", "data/ip2region.xdb")
	viper.SetDefault("geoip.ip2region.mode", "vector")
}

func initAPIKeyConfig() {
	apiKey := os.Getenv("SHORTENER_API_KEY")
	if apiKey == "" {
		apiKey = viper.GetString("server.api_key")
	}

	if apiKey == "" {
		apiKey = utils.GenerateCode(32)
	}

	shared.GlobalAPIKey = apiKey
}

func initUserConfig() {
	username := os.Getenv("SHORTENER_USERNAME")
	if username == "" {
		username = viper.GetString("admin.username")
	}

	password := os.Getenv("SHORTENER_PASSWORD")
	if password == "" {
		password = viper.GetString("admin.password")
	}

	// 如果是空值则自动生成随机账号和密码
	if username == "" {
		username = utils.GenerateCode(5)
	}

	if password == "" {
		password = utils.GenerateCode(10)
	}

	shared.GlobalUser = &types.User{
		Username: username,
		Password: password,
	}
}
