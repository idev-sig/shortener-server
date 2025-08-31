package shared

import (
	"sync"

	"gorm.io/gorm"

	"go.xoder.cn/shortener/internal/cache"
	"go.xoder.cn/shortener/internal/pkgs/geoip"
	"go.xoder.cn/shortener/internal/types"
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
