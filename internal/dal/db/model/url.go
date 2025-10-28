package model

// Url 短网址表
type Url struct {
	ID          int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                              // 主键ID
	ShortCode   string      `gorm:"column:short_code;type:varchar(16);uniqueIndex;not null" json:"short_code"` // 短码
	OriginalURL string      `gorm:"column:original_url;type:varchar(2048);not null" json:"original_url"`       // 原始URL
	Description string      `gorm:"column:description;type:varchar(255)" json:"description"`                   // 描述
	Status      int8        `gorm:"column:status;type:smallint;default:0;index;not null" json:"status"`        // 状态
	UpdatedAt   RFC3339Time `gorm:"column:updated_at;type:datetime;not null;index" json:"updated_at"`          // 更新时间
	CreatedAt   RFC3339Time `gorm:"column:created_at;type:datetime;not null;index" json:"created_at"`          // 创建时间
	Histories   []History   `gorm:"foreignKey:UrlID;constraint:OnDelete:CASCADE"`
}
