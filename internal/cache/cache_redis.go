package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"

	"go.xoder.cn/shortener/internal/types"
)

// RedisCache 缓存
type RedisCache struct {
	config *types.CfgCache
	client *redis.Client
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(cfg *types.CfgCache, redisCfg *types.CfgCacheRedis) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host + ":" + strconv.Itoa(redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	return &RedisCache{config: cfg, client: client}, client.Ping(context.Background()).Err()
}

// Ping 检查缓存连接
func (t *RedisCache) Ping() error {
	return t.client.Ping(context.Background()).Err()
}

// Set 设置缓存
func (t *RedisCache) Set(key string, value any, ttl ...time.Duration) error {
	jsonBytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}

	expire := 0
	if len(ttl) > 0 {
		expire = int(ttl[0])
	}
	return t.client.Set(context.Background(), key, string(jsonBytes), time.Duration(expire)*time.Second).Err()
}

// Get 获取缓存
func (t *RedisCache) Get(key string) (string, error) {
	data, err := t.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

// Delete 删除缓存
func (t *RedisCache) Delete(key string) error {
	return t.client.Del(context.Background(), key).Err()
}

// ClearPrefix 清理缓存前缀
func (t *RedisCache) ClearPrefix(prefix string) error {
	// log.Printf("clear prefix: %+v", shared.GlobalCache)
	ctx := context.Background()
	iter := t.client.Scan(ctx, 0, prefix+"*", 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		if err := t.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	return nil
}

// 批量设置缓存
func (t *RedisCache) BatchSet(values map[string]string, ttl ...time.Duration) error {
	ctx := context.Background()
	pipe := t.client.Pipeline()

	expire := 0
	if len(ttl) > 0 {
		expire = int(ttl[0])
	}

	for key, value := range values {
		// log.Printf("value: %v", value)
		pipe.Set(ctx, key, value, time.Duration(expire)*time.Second)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}
