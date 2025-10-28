package model

// History 短网址访问记录表
type History struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                       // 主键ID
	UrlID      int64       `gorm:"column:url_id;index;not null;index:idx_url_id" json:"url_id"`                        // 对应的短链接ID
	ShortCode  string      `gorm:"column:short_code;type:varchar(16);not null;index:idx_short_code" json:"short_code"` // 短码
	IPAddress  string      `gorm:"column:ip_address;type:varchar(45);not null;index:idx_ip_address" json:"ip_address"` // 访问者IP（IPv6最大45字符）
	UserAgent  string      `gorm:"column:user_agent;type:text;not null" json:"user_agent"`                             // 用户代理信息
	Referer    string      `gorm:"column:referer;type:text" json:"referer"`                                            // 来源URL
	Country    string      `gorm:"column:country;type:varchar(100)" json:"country"`                                    // 国家
	Region     string      `gorm:"column:region;type:varchar(100)" json:"region"`                                      // 地区/省份
	Province   string      `gorm:"column:province;type:varchar(100)" json:"province"`                                  // 省份
	City       string      `gorm:"column:city;type:varchar(100)" json:"city"`                                          // 城市
	ISP        string      `gorm:"column:isp;type:varchar(100)" json:"isp"`                                            // 运营商
	DeviceType string      `gorm:"column:device_type;type:varchar(50)" json:"device_type"`                             // 设备类型（pc/mobile/tablet）
	OS         string      `gorm:"column:os;type:varchar(50)" json:"os"`                                               // 操作系统
	Browser    string      `gorm:"column:browser;type:varchar(50)" json:"browser"`                                     // 浏览器类型
	AccessedAt RFC3339Time `gorm:"column:accessed_at;type:datetime;not null;index" json:"accessed_at"`                 // 访问时间
	CreatedAt  RFC3339Time `gorm:"column:created_at;type:datetime;not null;index" json:"created_at"`                   // 创建时间
	Url        Url         `gorm:"foreignKey:UrlID"`
}
