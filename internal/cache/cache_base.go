package cache

import (
	"strings"
	"sync"
	"time"

	"go.xoder.cn/shortener-server/internal/ecodes"
)

// 缓存项结构
type baseCacheItem struct {
	Value      any
	Expiration int64 // 过期时间（UnixNano 时间戳）
}

// 判断缓存项是否过期
func (item baseCacheItem) Expired() bool {
	if item.Expiration == 0 {
		return false // 永不过期
	}
	return time.Now().UnixNano() > item.Expiration
}

type BaseCache struct {
	items map[string]baseCacheItem
	mu    sync.RWMutex
}

func NewBaseCache() (*BaseCache, error) {
	c := &BaseCache{
		items: make(map[string]baseCacheItem),
	}
	go c.cleanupExpired()
	return c, nil
}

func (c *BaseCache) Items() map[string]baseCacheItem {
	return c.items
}

// 定期清理过期缓存（每5分钟执行一次）
func (c *BaseCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if item.Expired() {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

func (t *BaseCache) Ping() error {
	return nil
}

func (t *BaseCache) Set(key string, value any, ttl ...time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var expiration int64
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0]).UnixNano()
	}
	t.items[key] = baseCacheItem{
		Value:      value,
		Expiration: expiration,
	}
	return nil
}

func (t *BaseCache) Get(key string) (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	item, exists := t.items[key]
	if !exists || item.Expired() {
		return "", ecodes.GetGeneralError(ecodes.ErrCodeCacheKeyNotFound)
	}
	return item.Value.(string), nil
}

func (t *BaseCache) Delete(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.items, key)
	return nil
}

func (t *BaseCache) ClearPrefix(prefix string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for key := range t.items {
		if strings.HasPrefix(key, prefix) {
			delete(t.items, key)
		}
	}
	return nil
}

func (t *BaseCache) BatchSet(values map[string]string, ttl ...time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var expiration int64
	if len(ttl) > 0 {
		expiration = time.Now().Add(ttl[0]).UnixNano()
	}
	for key, value := range values {
		t.items[key] = baseCacheItem{
			Value:      value,
			Expiration: expiration,
		}
	}
	return nil
}
