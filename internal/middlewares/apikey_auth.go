package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener/internal/ecodes"
	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/types"
)

// ApiKeyAuth 检查请求头中的 API Key
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if shared.GlobalAPIKey != "" && apiKey != shared.GlobalAPIKey {
			errCode := ecodes.ErrCodeUnauthorized
			c.JSON(http.StatusUnauthorized, types.ResErr{
				ErrorCode:    fmt.Sprintf("%d", errCode),
				ErrorMessage: ecodes.GetErrCodeMessage(errCode),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
