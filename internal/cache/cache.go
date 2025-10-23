package cache

import (
	"time"

	"go.xoder.cn/shortener-server/internal/ecodes"
)

// Cache 缓存
type Cache interface {
	Ping() error
	Set(key string, value any, ttl ...time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	ClearPrefix(prefix string) error
	BatchSet(values map[string]string, ttl ...time.Duration) error
}

// CacheManager 缓存管理器
type CacheManager struct {
	Enabled bool
	Cache   Cache
	Prefix  string
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(enabled bool, cache Cache, prefix string) *CacheManager {
	return &CacheManager{Enabled: enabled, Cache: cache, Prefix: prefix}
}

// Get 获取缓存
func (c *CacheManager) Get(key string) (string, error) {
	if !c.Enabled {
		return "", ecodes.ErrCacheDisabled
	}
	return c.Cache.Get(key)
}

// Set 设置缓存
func (c *CacheManager) Set(key string, value any, ttl ...time.Duration) error {
	if !c.Enabled {
		return ecodes.ErrCacheDisabled
	}
	return c.Cache.Set(key, value, ttl...)
}

// Delete 删除缓存
func (c *CacheManager) Delete(key string) error {
	if !c.Enabled {
		return ecodes.ErrCacheDisabled
	}
	return c.Cache.Delete(key)
}

// ClearPrefix 清理缓存前缀
func (c *CacheManager) ClearPrefix(prefix string) error {
	if !c.Enabled {
		return ecodes.ErrCacheDisabled
	}
	return c.Cache.ClearPrefix(prefix)
}

// BatchSet 批量设置缓存
func (c *CacheManager) BatchSet(values map[string]string, ttl ...time.Duration) error {
	if !c.Enabled {
		return ecodes.ErrCacheDisabled
	}
	return c.Cache.BatchSet(values, ttl...)
}

// Ping 检查缓存连接
func (c *CacheManager) Ping() error {
	if !c.Enabled {
		return ecodes.ErrCacheDisabled
	}
	return c.Cache.Ping()
}

// GetKey 获取缓存key
func (c *CacheManager) GetKey(key string) string {
	if !c.Enabled {
		return ""
	}
	return c.Prefix + key
}
