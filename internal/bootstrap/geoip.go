package bootstrap

import (
	"fmt"

	"github.com/spf13/viper"

	"go.bdev.cn/shortener/internal/pkgs/geoip"
	"go.bdev.cn/shortener/internal/shared"
)

// initGeoIP 初始化IP地址库
func initGeoIP() {
	enabled := viper.GetBool("geoip.enabled")
	if !enabled {
		return
	}

	var geoIP geoip.GeoIP
	geoIPType := viper.GetString("geoip.type")
	switch geoIPType {
	case "ip2region":
		geoIP = ip2RegionGeoIP()
	default:
		panic(fmt.Sprintf("geoip type not support: %s\n", geoIPType))
	}

	shared.GlobalGeoIP = geoip.NewGeoIPManager(enabled, geoIPType, geoIP)
}

// initIP2Region 初始化IP2Region
func ip2RegionGeoIP() *geoip.IP2Region {
	dbPath := viper.GetString("geoip.ip2region.path")
	loadMode := viper.GetString("geoip.ip2region.mode")

	ip2region, err := geoip.NewIP2Region(dbPath, loadMode)
	if err != nil {
		panic(fmt.Sprintf("failed to create ip2region: %s\n", err))
	}
	return ip2region
}
