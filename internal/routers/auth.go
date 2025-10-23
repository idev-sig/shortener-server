package routers

import (
	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener-server/internal/middlewares"
	"go.xoder.cn/shortener-server/internal/shared"
)

func authMiddleware() gin.HandlerFunc {
	apiKeyAuth := &middlewares.APIKeyAuth{
		ValidKeys: map[string]bool{shared.GlobalAPIKey: true},
		Header:    "X-API-KEY",
		Query:     "api_key",
	}

	return middlewares.MultiAuthMiddleware(apiKeyAuth, &middlewares.BearerTokenAuth{})
}
