package shared

import (
	"sync"

	"gorm.io/gorm"

	"go.bdev.cn/shortener/internal/cache"
	"go.bdev.cn/shortener/internal/pkgs/geoip"
	"go.bdev.cn/shortener/internal/types"
)

var (
	GlobalShorten *types.CfgShorten
	GlobalDB      *gorm.DB
	GlobalAPIKey  string
	GlobalCache   *cache.CacheManager
	GlobalGeoIP   *geoip.GeoIPManager

	GlobalUser      *types.User
	GlobalUserCache sync.Map
)
