package logics

import (
	"fmt"
	"time"

	"github.com/ua-parser/uap-go/uaparser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"go.xoder.cn/shortener/internal/dal/db/model"
	"go.xoder.cn/shortener/internal/ecodes"
	"go.xoder.cn/shortener/internal/pkgs/geoip"
	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/types"
	"go.xoder.cn/shortener/internal/utils"
)

// HistoryLogic 历史记录逻辑层
type HistoryLogic struct {
	logic
	geoip *geoip.GeoIPManager
}

// NewHistoryLogic 创建历史记录逻辑层
func NewHistoryLogic() *HistoryLogic {
	t := &HistoryLogic{
		geoip: shared.GlobalGeoIP,
	}
	t.db = shared.GlobalDB
	return t
}

// HistoryAdd 添加历史记录
func (t *HistoryLogic) HistoryAdd(params types.HistoryParams) error {
	nowTime := time.Now().Local()

	// 解析用户代理
	uaParser := uaparser.NewFromSaved()
	client := uaParser.Parse(params.UserAgent)
	deviceType := cases.Title(language.English).String(
		simplifyDeviceType(client.Device.ToString()),
	)

	// 初始化地理位置信息
	geoInfo := struct {
		Country  string
		Region   string
		Province string
		City     string
		ISP      string
	}{}

	if t.geoip != nil && t.geoip.Enabled {
		if ipByte, err := t.geoip.IPStr2Byte(params.IPAddress); err == nil {
			if ipInfo, err := t.geoip.Search(ipByte); err == nil {
				ipData := t.geoip.Parse(ipInfo)
				geoInfo.Country = ipData.Country
				geoInfo.Region = ipData.Region
				geoInfo.Province = ipData.Province
				geoInfo.City = ipData.City
				geoInfo.ISP = ipData.ISP
			}
		}
	}

	history := model.History{
		UrlID:      params.URLID,
		ShortCode:  params.ShortCode,
		IPAddress:  params.IPAddress,
		UserAgent:  params.UserAgent,
		Referer:    params.Referer,
		Country:    geoInfo.Country,
		Region:     geoInfo.Region,
		Province:   geoInfo.Province,
		City:       geoInfo.City,
		ISP:        geoInfo.ISP,
		DeviceType: deviceType,
		OS:         client.Os.ToString(),
		Browser:    client.UserAgent.ToString(),
		AccessedAt: nowTime,
		CreatedAt:  nowTime,
	}
	// log.Printf("history: %+v\n", history)

	return t.db.Create(&history).Error
}

// HistoryDeleteAll 删除所有历史记录
func (t *HistoryLogic) HistoryDeleteAll(ids []string) int {
	if res := t.db.Where("id in (?)", ids).Delete(&model.History{}); res.Error != nil {
		return ecodes.ErrCodeDatabaseError
	}

	return ecodes.ErrCodeSuccess
}

// HistoryAll 获取所有短链接
func (t *HistoryLogic) HistoryAll(reqQuery types.ReqQueryHistory) (int, []types.ResHistory, types.ResPage) {
	results := make([]types.ResHistory, 0)
	pageInfo := types.ResPage{}

	// 查询数据库
	query := t.db.Model(&model.History{}).
		Order(fmt.Sprintf("%s %s", reqQuery.SortBy, reqQuery.Order))

	if reqQuery.Code != "" {
		query = query.Where("short_code = ?", reqQuery.Code)
	}

	if reqQuery.IP != "" {
		query = query.Where("ip_address = ?", reqQuery.IP)
	}

	// 计算总条数
	var total int64
	query = query.Count(&total)
	if query.Error != nil {
		return ecodes.ErrCodeDatabaseError, results, pageInfo
	}

	// 分页查询
	data := make([]model.History, 0)
	resDB := query.Offset(int((reqQuery.Page - 1) * reqQuery.PageSize)).
		Limit(int(reqQuery.PageSize)).
		Find(&data)
	if resDB.Error != nil {
		return ecodes.ErrCodeDatabaseError, results, pageInfo
	}

	// 页码信息
	pageInfo.Page = reqQuery.Page
	pageInfo.PageSize = reqQuery.PageSize
	pageInfo.CurrentCount = resDB.RowsAffected
	pageInfo.TotalItems = total
	pageInfo.TotalPages = total / int64(reqQuery.PageSize)
	if total%int64(reqQuery.PageSize) != 0 {
		pageInfo.TotalPages++
	}

	for _, item := range data {
		results = append(results, types.ResHistory{
			ID:         item.ID,
			UrlID:      item.UrlID,
			ShortCode:  item.ShortCode,
			IPAddress:  item.IPAddress,
			UserAgent:  item.UserAgent,
			Referer:    item.Referer,
			Country:    item.Country,
			Region:     item.Region,
			Province:   item.Province,
			City:       item.City,
			ISP:        item.ISP,
			DeviceType: item.DeviceType,
			OS:         item.OS,
			Browser:    item.Browser,
			AccessedAt: utils.TimeToStr(item.AccessedAt),
			CreatedAt:  utils.TimeToStr(item.CreatedAt),
		})
	}

	return ecodes.ErrCodeSuccess, results, pageInfo
}
