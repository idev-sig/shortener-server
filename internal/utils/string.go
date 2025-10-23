package utils

import (
	"fmt"
	"math/rand"

	"go.xoder.cn/shortener-server/internal/shared"
)

// GenerateCode 生成短码(6位)
func GenerateCode(length int) string {
	charset := shared.GlobalShorten.Charset

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		result[i] = charset[index]
	}
	return string(result)
}

// IntToString 将整数转换为字符串
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}
