package utils

import (
	"math/rand"

	"go.xoder.cn/shortener/internal/shared"
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
