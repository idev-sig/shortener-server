package types

// ReqCode URL Path
type ReqCode struct {
	Code string `uri:"code" binding:"required"`
}

// ReqQuery 请求参数结构体
type ReqQuery struct {
	Page    int64  `form:"page,default=1" binding:"min=1"`
	PerPage int64  `form:"per_page,default=10" binding:"min=1,max=100"`
	SortBy  string `form:"sort_by,omitempty"`
	Order   string `form:"order,omitempty" binding:"omitempty,oneof=asc desc"`
}

type ReqQueryShorten struct {
	ReqQuery
	ShortCode   string `form:"short_code,omitempty" binding:"omitempty"`
	OriginalURL string `form:"original_url,omitempty" binding:"omitempty"`
	Status      int64  `form:"status,omitempty,default=-1" binding:"omitempty"`
}

type ReqQueryHistory struct {
	ReqQuery
	ShortCode string `form:"short_code,omitempty" binding:"omitempty"`
	IPAddress string `form:"ip_address,omitempty" binding:"omitempty"`
}

// ResShorten 短链接响应
type ResShorten struct {
	ID          int64  `json:"id"`
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	Description string `json:"description,omitempty"`
	Status      int8   `json:"status"` // 0=启用, 1=禁用, 2=未知
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ResHistory 历史记录响应
type ResHistory struct {
	ID         int64  `json:"id"`
	UrlID      int64  `json:"url_id"`
	ShortCode  string `json:"short_code"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
	Referer    string `json:"referer,omitempty"`
	Country    string `json:"country,omitempty"`
	Region     string `json:"region,omitempty"`
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
	ISP        string `json:"isp,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	OS         string `json:"os,omitempty"`
	Browser    string `json:"browser,omitempty"`
	AccessedAt string `json:"accessed_at"`
	CreatedAt  string `json:"created_at"`
}

// ResPage 分页响应 (对应 OpenAPI PageMeta)
type ResPage struct {
	Page       int64 `json:"page"`        // 当前页码
	PerPage    int64 `json:"per_page"`    // 每页数量
	Count      int64 `json:"count"`       // 当前页条目数
	Total      int64 `json:"total"`       // 总条目数
	TotalPages int64 `json:"total_pages"` // 总页数
}

// ResErr 错误响应 (对应 OpenAPI ErrorResponse)
type ResErr struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ResSuccess 成功响应
type ResSuccess[T any] struct {
	Data T       `json:"data"` // 数据
	Meta ResPage `json:"meta"` // 元数据
}

// CfgShorten 短链接配置
type CfgShorten struct {
	Length  int    `json:"length"`
	Charset string `json:"charset"`
}

// CfgCache 缓存配置
type CfgCache struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Expire  int    `json:"expire"`
	Prefix  string `json:"prefix"`
}

// CfgCacheRedis 缓存配置
type CfgCacheRedis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// CfgCacheValkey 缓存配置
type CfgCacheValkey struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
