package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/valkey-io/valkey-go"

	"go.bdev.cn/shortener/internal/types"
)

// ValkeyCache 缓存
type ValkeyCache struct {
	config      *types.CfgCache
	clientPoint *valkey.Client
	client      valkey.Client
}

// NewValkeyCache 创建Valkey缓存
func NewValkeyCache(cfg *types.CfgCache, valkeyCfg *types.CfgCacheValkey) (*ValkeyCache, error) {
	address := valkeyCfg.Host + ":" + strconv.Itoa(valkeyCfg.Port)
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{address},
		Username:    valkeyCfg.Username,
		Password:    valkeyCfg.Password,
		SelectDB:    valkeyCfg.DB,
	})
	if err != nil {
		return nil, err
	}
	return &ValkeyCache{config: cfg, clientPoint: &client, client: client}, nil
}

// Ping 检查缓存连接
func (t *ValkeyCache) Ping() error {
	return t.client.Do(context.Background(), t.client.B().Ping().Build()).Error()
}

// Set 设置缓存
func (t *ValkeyCache) Set(key string, value any, ttl ...time.Duration) error {
	jsonBytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}

	expire := 0
	if len(ttl) > 0 {
		expire = int(ttl[0])
	}
	ctx := context.Background()
	builder := t.client.B().Set().Key(key).Value(string(jsonBytes))
	if expire > 0 {
		builder.ExSeconds(int64(expire))
	}
	return t.client.Do(ctx, builder.Build()).Error()
}

// Get 获取缓存
func (t *ValkeyCache) Get(key string) (string, error) {
	resp, err := t.client.Do(context.Background(), t.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "", err
	}
	return resp, nil
}

// Delete 删除缓存
func (t *ValkeyCache) Delete(key string) error {
	return t.client.Do(context.Background(), t.client.B().Del().Key(key).Build()).Error()
}

// ClearPrefix 清理缓存前缀
func (t *ValkeyCache) ClearPrefix(prefix string) error {
	// log.Printf("clear prefix: %+v", shared.GlobalCache)
	var cursor uint64 = 0
	prefix = prefix + "*"
	ctx := context.Background()
	keys := make([]string, 0)

	for {
		scanCmd := t.client.B().Scan().
			Cursor(cursor).
			Match(prefix).
			Count(1000).
			Build()
		resp := t.client.Do(ctx, scanCmd)
		se, err := resp.AsScanEntry()
		if err != nil {
			return err
		}
		keys = append(keys, se.Elements...)
		cursor = se.Cursor
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		t.client.Do(ctx, t.client.B().Del().Key(keys...).Build())
	}

	return nil
}

// 批量设置缓存
func (t *ValkeyCache) BatchSet(values map[string]string, ttl ...time.Duration) error {
	ctx := context.Background()

	expire := 0
	if len(ttl) > 0 {
		expire = int(ttl[0])
	}

	cmds := make(valkey.Commands, 0, len(values))
	for key, value := range values {
		builder := t.client.B().Set().
			Key(key).
			Value(value)

		if expire > 0 {
			builder.Ex(time.Duration(expire) * time.Second)
		}
		cmds = append(cmds, builder.Build())
	}

	responses := t.client.DoMulti(ctx, cmds...)
	for _, resp := range responses {
		if err := resp.Error(); err != nil {
			return err
		}
	}

	return nil
}
