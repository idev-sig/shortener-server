package model

import (
	"time"
)

// Url 短网址表
type Url struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                 // 主键ID
	ShortCode   string    `gorm:"column:short_code;type:varchar(16);uniqueIndex;not null" json:"short_code"`    // 短码
	OriginalURL string    `gorm:"column:original_url;type:varchar(2048);not null" json:"original_url"`          // 原始URL
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`                      // 描述
	Status      int8      `gorm:"column:status;type:smallint;default:0;index;not null" json:"status"`           // 状态
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;precision:6;not null;index" json:"updated_at"` // 更新时间
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;precision:6;not null;index" json:"created_at"` // 创建时间
	Histories   []History `gorm:"foreignKey:UrlID;constraint:OnDelete:CASCADE"`
}

// // 按需添加以下索引
// func (Url) Indexes() []string {
// 	return []string{
// 		"idx_short_code",
// 		"idx_status",
// 		"idx_created_at",
// 		"idx_updated_at",
// 		"CREATE INDEX idx_country_device ON request_histories(country_code,device_type)",
// 	}
// }
