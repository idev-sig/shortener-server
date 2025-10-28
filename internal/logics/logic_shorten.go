package logics

import (
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"

	"go.xoder.cn/shortener-server/internal/dal/db/model"
	"go.xoder.cn/shortener-server/internal/ecodes"
	"go.xoder.cn/shortener-server/internal/types"
	"go.xoder.cn/shortener-server/internal/utils"
)

// ShortenLogic 短链接逻辑层
type ShortenLogic struct {
	logic
}

// NewShortenLogic 创建短链接逻辑层
func NewShortenLogic() *ShortenLogic {
	t := &ShortenLogic{}
	t.init()
	return t
}

// ShortenAdd 添加短链接
func (t *ShortenLogic) ShortenAdd(code string, originalURL string, description string) (int, types.ResShorten) {
	result := types.ResShorten{}
	existingURL := model.Url{}

	// 1. 检查短码是否已存在（使用 GORM 的 Find 直接判断）
	if err := t.db.Where("short_code = ?", code).First(&existingURL).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ecodes.ErrCodeDatabaseError, result // 数据库查询错误
		}
		// 短码不存在，继续流程
	} else {
		return ecodes.ErrCodeConflict, result // 短码已存在
	}

	// 2. 创建新记录
	nowTime := model.RFC3339Time{Time: time.Now().UTC()}
	newURL := model.Url{
		ShortCode:   code,
		OriginalURL: originalURL,
		Description: description,
		Status:      0,
		CreatedAt:   nowTime,
		UpdatedAt:   nowTime,
	}

	if err := t.db.Create(&newURL).Error; err != nil {
		return ecodes.ErrCodeDatabaseError, result // 创建失败
	}

	// 3. 缓存短链接
	if err := t.cache.Set(t.cache.GetKey(newURL.ShortCode), newURL); err != nil && !errors.Is(err, ecodes.ErrCacheDisabled) {
		return ecodes.ErrCodeCacheError, result // 缓存失败
	}

	// 4. 构造返回结果
	result = types.ResShorten{
		ID:          newURL.ID,
		ShortCode:   newURL.ShortCode,
		ShortURL:    t.GetSiteURL(newURL.ShortCode),
		OriginalURL: newURL.OriginalURL,
		Description: newURL.Description,
		Status:      newURL.Status,
		CreatedAt:   utils.TimeToStr(nowTime.Time),
		UpdatedAt:   utils.TimeToStr(nowTime.Time),
	}

	return ecodes.ErrCodeSuccess, result
}

// ShortenDelete 删除短链接
func (t *ShortenLogic) ShortenDelete(code string) int {
	if res := t.db.Where("short_code = ?", code).Delete(&model.Url{}); res.Error != nil {
		return ecodes.ErrCodeDatabaseError
	} else if res.RowsAffected == 0 {
		return ecodes.ErrCodeNotFound
	}

	// 删除缓存
	if err := t.cache.Delete(t.cache.GetKey(code)); err != nil && !errors.Is(err, ecodes.ErrCacheDisabled) {
		return ecodes.ErrCodeCacheError // 缓存删除失败
	}

	return ecodes.ErrCodeSuccess
}

// ShortenDeleteAll 删除所有短链接
func (t *ShortenLogic) ShortenDeleteAll(ids []string) int {
	if res := t.db.Where("id in (?)", ids).Delete(&model.Url{}); res.Error != nil {
		return ecodes.ErrCodeDatabaseError
	}

	// 删除缓存
	for _, id := range ids {
		if err := t.cache.Delete(t.cache.GetKey(id)); err != nil && !errors.Is(err, ecodes.ErrCacheDisabled) {
			return ecodes.ErrCodeCacheError // 缓存删除失败
		}
	}

	return ecodes.ErrCodeSuccess
}

// ShortenUpdate 更新短链接
func (t *ShortenLogic) ShortenUpdate(code string, originalURL string, description string) (int, types.ResShorten) {
	result := types.ResShorten{}

	var existingURL model.Url
	if err := t.db.Where("short_code = ?", code).First(&existingURL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ecodes.ErrCodeNotFound, result
		}
		return ecodes.ErrCodeDatabaseError, result
	}

	// 准备更新字段
	updates := make(map[string]any)

	if originalURL != "" {
		updates["original_url"] = originalURL
	}
	if description != "" {
		updates["description"] = description
	}

	nowTime := model.RFC3339Time{Time: time.Now().UTC()}
	updates["updated_at"] = nowTime

	if err := t.db.Model(&existingURL).Updates(updates).Error; err != nil {
		return ecodes.ErrCodeDatabaseError, result
	}

	if err := t.cache.Set(t.cache.GetKey(existingURL.ShortCode), existingURL); err != nil && !errors.Is(err, ecodes.ErrCacheDisabled) {
		return ecodes.ErrCodeCacheError, result // 缓存失败
	}

	result = types.ResShorten{
		ID:          existingURL.ID,
		ShortCode:   existingURL.ShortCode,
		ShortURL:    t.GetSiteURL(existingURL.ShortCode),
		OriginalURL: existingURL.OriginalURL,
		Description: existingURL.Description,
		Status:      existingURL.Status,
		UpdatedAt:   utils.TimeToStr(nowTime.Time),
		CreatedAt:   utils.TimeToStr(existingURL.CreatedAt.Time),
	}

	return ecodes.ErrCodeSuccess, result
}

// ShortenFind 获取短链接
func (t *ShortenLogic) ShortenFind(code string) (int, types.ResShorten) {
	var data model.Url

	// 1. 从缓存中获取
	cacheKey := t.cache.GetKey(code)
	if cacheData, err := t.cache.Get(cacheKey); err == nil {
		// log.Printf("cacheData: %v", cacheData)
		if err := sonic.Unmarshal([]byte(cacheData), &data); err != nil {
			return ecodes.ErrCodeCacheError, types.ResShorten{} // 缓存反序列化失败
		}
	} else {
		// 从数据库中获取
		if err := t.db.Where("short_code = ?", code).First(&data).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ecodes.ErrCodeNotFound, types.ResShorten{}
			}
			return ecodes.ErrCodeDatabaseError, types.ResShorten{}
		}

		// log.Printf("data: %v", data)
		// 缓存短链接
		if err := t.cache.Set(cacheKey, data); err != nil && !errors.Is(err, ecodes.ErrCacheDisabled) {
			return ecodes.ErrCodeCacheError, types.ResShorten{} // 缓存失败
		}
	}

	result := types.ResShorten{
		ID:          data.ID,
		ShortCode:   data.ShortCode,
		ShortURL:    t.GetSiteURL(data.ShortCode),
		OriginalURL: data.OriginalURL,
		Description: data.Description,
		Status:      data.Status,
		CreatedAt:   utils.TimeToStr(data.CreatedAt.Time),
		UpdatedAt:   utils.TimeToStr(data.UpdatedAt.Time),
	}

	return ecodes.ErrCodeSuccess, result
}

// ShortenAll 获取所有短链接
func (t *ShortenLogic) ShortenAll(reqQuery types.ReqQueryShorten) (int, []types.ResShorten, types.ResPage) {
	results := make([]types.ResShorten, 0)
	pageInfo := types.ResPage{}

	// 查询数据库
	query := t.db.Model(&model.Url{}).
		Order(fmt.Sprintf("%s %s", reqQuery.SortBy, reqQuery.Order))

	if reqQuery.ShortCode != "" {
		query = query.Where("short_code = ?", reqQuery.ShortCode)
	}

	if reqQuery.OriginalURL != "" {
		// query = query.Where("original_url = ?", reqQuery.OriginalURL)
		// 模糊查找
		query = query.Where("original_url like ?", "%"+reqQuery.OriginalURL+"%")
	}

	if reqQuery.Status != -1 {
		query = query.Where("status = ?", reqQuery.Status)
	}

	// 计算总条数
	var total int64
	query = query.Count(&total)
	if query.Error != nil {
		return ecodes.ErrCodeDatabaseError, results, pageInfo
	}

	// 分页查询
	data := make([]model.Url, 0)
	resDB := query.Offset(int((reqQuery.Page - 1) * reqQuery.PerPage)).
		Limit(int(reqQuery.PerPage)).
		Find(&data)
	if resDB.Error != nil {
		return ecodes.ErrCodeDatabaseError, results, pageInfo
	}

	// 页码信息
	pageInfo.Page = reqQuery.Page
	pageInfo.PerPage = reqQuery.PerPage
	pageInfo.Count = resDB.RowsAffected
	pageInfo.Total = total
	pageInfo.TotalPages = total / int64(reqQuery.PerPage)
	if total%int64(reqQuery.PerPage) != 0 {
		pageInfo.TotalPages++
	}

	for _, item := range data {
		results = append(results, types.ResShorten{
			ID:          item.ID,
			ShortCode:   item.ShortCode,
			ShortURL:    t.GetSiteURL(item.ShortCode),
			OriginalURL: item.OriginalURL,
			Description: item.Description,
			Status:      item.Status,
			CreatedAt:   utils.TimeToStr(item.CreatedAt.Time),
			UpdatedAt:   utils.TimeToStr(item.UpdatedAt.Time),
		})
	}

	return ecodes.ErrCodeSuccess, results, pageInfo
}
