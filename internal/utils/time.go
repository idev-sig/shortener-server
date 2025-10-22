package utils

import "time"

// TimestampToTime 时间戳转时间
func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// TimeToStr 时间转字符串 (ISO 8601 格式)
func TimeToStr(nowTime time.Time) string {
	return nowTime.UTC().Format(time.RFC3339)
}

// StrToTime 字符串转时间 (ISO 8601 格式)
func StrToTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

// StrToTimestamp 字符串转时间戳
func StrToTimestamp(str string) (int64, error) {
	time, err := StrToTime(str)
	return time.Unix(), err
}

// TimestampToStr 时间戳转字符串 (ISO 8601 格式)
func TimestampToStr(timestamp int64) string {
	return TimestampToTime(timestamp).UTC().Format(time.RFC3339)
}
