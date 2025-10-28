package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// RFC3339Time 自定义时间类型，以 RFC3339 格式存储到数据库
type RFC3339Time struct {
	time.Time
}

// Scan 实现 sql.Scanner 接口，从数据库读取时间
func (t *RFC3339Time) Scan(value any) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v.UTC()
		return nil
	case string:
		// 尝试解析 RFC3339 格式
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			// 如果失败，尝试其他常见格式
			parsed, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
			if err != nil {
				parsed, err = time.Parse("2006-01-02 15:04:05", v)
				if err != nil {
					return fmt.Errorf("failed to parse time: %v", err)
				}
			}
		}
		t.Time = parsed.UTC()
		return nil
	case []byte:
		return t.Scan(string(v))
	default:
		return fmt.Errorf("unsupported type for RFC3339Time: %T", value)
	}
}

// Value 实现 driver.Valuer 接口，写入数据库时使用 RFC3339 格式
func (t RFC3339Time) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	// 返回 RFC3339 格式的字符串
	return t.Time.UTC().Format(time.RFC3339), nil
}

// MarshalJSON 实现 JSON 序列化
func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.UTC().Format(time.RFC3339) + `"`), nil
}

// UnmarshalJSON 实现 JSON 反序列化
func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	// 移除引号
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	t.Time = parsed.UTC()
	return nil
}
