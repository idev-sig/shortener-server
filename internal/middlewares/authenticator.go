package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener/internal/ecodes"
	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/types"
)

type Authenticator interface {
	Authenticate(c *gin.Context) (bool, error)
}

type APIKeyAuth struct {
	ValidKeys map[string]bool
	Header    string
	Query     string
}

func (a *APIKeyAuth) Authenticate(c *gin.Context) (bool, error) {
	if key := c.GetHeader(a.Header); key != "" && a.ValidKeys[key] {
		return true, nil
	}
	if key := c.Query(a.Query); key != "" && a.ValidKeys[key] {
		return true, nil
	}
	return false, nil
}

type BearerTokenAuth struct{}

func (b *BearerTokenAuth) Authenticate(c *gin.Context) (bool, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false, nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader { // 没有Bearer前缀
		return false, nil
	}

	userInfo, ok := shared.GlobalUserCache.Load(token)
	if !ok || userInfo != shared.GlobalUser.Username {
		return false, nil
	}

	return true, nil
}

func MultiAuthMiddleware(authenticators ...Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, auth := range authenticators {
			success, err := auth.Authenticate(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, types.ResErr{
					ErrCode: ecodes.ErrCodeUserAuthFailed,
					ErrInfo: ecodes.GetErrCodeMessage(ecodes.ErrCodeUserAuthFailed),
				})
				return
			}
			if success {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, types.ResErr{
			ErrCode: ecodes.ErrCodeUnauthorized,
			ErrInfo: ecodes.GetErrCodeMessage(ecodes.ErrCodeUnauthorized),
		})
	}
}
