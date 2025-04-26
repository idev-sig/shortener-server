package routers

import (
	"github.com/gin-gonic/gin"

	"go.dsig.cn/shortener/internal/middlewares"
	"go.dsig.cn/shortener/internal/shared"
)

func authMiddleware() gin.HandlerFunc {
	apiKeyAuth := &middlewares.APIKeyAuth{
		ValidKeys: map[string]bool{shared.GlobalAPIKey: true},
		Header:    "X-API-KEY",
		Query:     "api_key",
	}

	return middlewares.MultiAuthMiddleware(apiKeyAuth, &middlewares.BearerTokenAuth{})
}
