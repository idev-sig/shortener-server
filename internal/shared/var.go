package shared

import (
	"sync"

	"gorm.io/gorm"

	"go.xoder.cn/shortener-server/internal/cache"
	"go.xoder.cn/shortener-server/internal/pkgs/geoip"
	"go.xoder.cn/shortener-server/internal/types"
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
