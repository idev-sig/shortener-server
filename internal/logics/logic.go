package logics

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"go.xoder.cn/shortener/internal/cache"
	"go.xoder.cn/shortener/internal/shared"
)

// logic 逻辑层
type logic struct {
	db       *gorm.DB
	cache    *cache.CacheManager
	site_url string
}

// init 初始化
func (t *logic) init() {
	t.db = shared.GlobalDB
	t.cache = shared.GlobalCache
	t.site_url = viper.GetString("server.site_url")
	if t.site_url == "" {
		t.site_url = "http://" + viper.GetString("server.address")
	}
}

// GetSiteURL 获取短链接的完整URL
func (t *logic) GetSiteURL(code string) string {
	return t.site_url + "/" + code
}
