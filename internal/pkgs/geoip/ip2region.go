package geoip

import (
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// IP2Region 结构体
type IP2Region struct {
	path     string
	mode     string
	searcher *xdb.Searcher
}

// NewIP2Region 创建 IP2Region 结构体
func NewIP2Region(dbPath string, loadMode string, version string) (*IP2Region, error) {
	var err error
	var searcher *xdb.Searcher

	xdbVersion := xdb.IPv4
	if version == "6" {
		xdbVersion = xdb.IPv6
	}

	switch loadMode {
	case "vector":
		searcher, err = loadVectorIndex(dbPath, xdbVersion)
	case "memory":
		searcher, err = loadMemoryIndex(dbPath, xdbVersion)
	case "file":
		searcher, err = loadFileIndex(dbPath, xdbVersion)
	default:
		return nil, fmt.Errorf("invalid load mode: %s", loadMode)
	}

	return &IP2Region{
		path:     dbPath,
		mode:     loadMode,
		searcher: searcher,
	}, err
}

// loadVectorIndex 加载向量索引
func loadVectorIndex(dbPath string, version *xdb.Version) (*xdb.Searcher, error) {
	vIndex, err := xdb.LoadVectorIndexFromFile(dbPath)
	if err != nil {
		return nil, err
	}
	return xdb.NewWithVectorIndex(version, dbPath, vIndex)
}

// loadMemoryIndex 加载内存索引
func loadMemoryIndex(dbPath string, version *xdb.Version) (*xdb.Searcher, error) {
	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		return nil, err
	}
	return xdb.NewWithBuffer(version, cBuff)
}

// loadFileIndex 加载文件索引
func loadFileIndex(dbPath string, version *xdb.Version) (*xdb.Searcher, error) {
	return xdb.NewWithFileOnly(version, dbPath)
}

// Search 搜索IP
func (t *IP2Region) Search(ip []byte) (string, error) {
	return t.searcher.Search(ip)
}

// SearchByStr 搜索IP
func (t *IP2Region) SearchByStr(ip string) (string, error) {
	return t.searcher.SearchByStr(ip)
}

// Parse 解析IP数据
func (t *IP2Region) Parse(data string) *GeoIPData {
	parts := strings.Split(data, "|")
	var country, province, city, isp string
	if len(parts) > 0 {
		country = parts[0]
	}
	if len(parts) > 1 {
		province = parts[1]
	}
	if len(parts) > 2 {
		city = parts[2]
	}
	if len(parts) > 3 {
		isp = parts[3]
	}

	return &GeoIPData{
		Country:  country,
		Province: province,
		City:     city,
		ISP:      isp,
	}
}
