package bootstrap

import (
	"github.com/bytedance/sonic"
	"github.com/spf13/viper"

	"go.xoder.cn/shortener/internal/cache"
	"go.xoder.cn/shortener/internal/dal/db/model"
	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/types"
)

var cacheCfg *types.CfgCache

// initCache 初始化缓存
func initCache() {
	if err := viper.UnmarshalKey("cache", &cacheCfg); err != nil {
		panic("cache config unmarshal failed: " + err.Error())
	}

	var cacheClient cache.Cache
	if cacheCfg.Enabled {
		switch cacheCfg.Type {
		case "redis":
			cacheClient = redisCache()
		case "valkey":
			cacheClient = valkeyCache()
		default:
			cacheCfg.Enabled = false
		}
	}

	// log.Printf("cache client: %+v cacheCfg: %+v", cacheClient, cacheCfg)
	shared.GlobalCache = cache.NewCacheManager(cacheCfg.Enabled, cacheClient, cacheCfg.Prefix)
	if !cacheCfg.Enabled {
		return
	}
	if err := shared.GlobalCache.Ping(); err != nil {
		panic("cache ping failed: " + err.Error())
	}

	loadAllShorten()
}

// redisCache 初始化 redis
func redisCache() *cache.RedisCache {
	var redisCfg types.CfgCacheRedis
	if err := viper.UnmarshalKey("cache.redis", &redisCfg); err != nil {
		panic("cache redis config unmarshal failed: " + err.Error())
	}
	cacheClient, err := cache.NewRedisCache(cacheCfg, &redisCfg)
	if err != nil {
		panic("cache redis init failed: " + err.Error())
	}
	return cacheClient
}

// valkeyCache 初始化 valkey
func valkeyCache() *cache.ValkeyCache {
	var valkeyCfg types.CfgCacheValkey
	if err := viper.UnmarshalKey("cache.valkey", &valkeyCfg); err != nil {
		panic("cache valkey config unmarshal failed: " + err.Error())
	}
	cacheClient, err := cache.NewValkeyCache(cacheCfg, &valkeyCfg)
	if err != nil {
		panic("cache valkey init failed: " + err.Error())
	}
	return cacheClient
}

// loadAllShorten 加载所有短链接
func loadAllShorten() {
	var shortens []model.Url
	if err := shared.GlobalDB.Find(&shortens).Error; err != nil {
		panic("load all shorten failed: " + err.Error())
	}

	if err := shared.GlobalCache.ClearPrefix(shared.GlobalCache.GetKey("")); err != nil {
		panic("cache clear prefix failed: " + err.Error())
	}

	items := make(map[string]string, len(shortens))
	for _, shorten := range shortens {
		item, _ := sonic.Marshal(shorten)
		items[shared.GlobalCache.GetKey(shorten.ShortCode)] = string(item)
	}

	if err := shared.GlobalCache.BatchSet(items); err != nil {
		panic("cache batch set failed: " + err.Error())
	}
}
