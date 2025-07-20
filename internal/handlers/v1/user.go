package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.bdev.cn/shortener/internal/shared"
)

// UserHandler 用户处理器
type UserHandler struct {
	handler
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	t := &UserHandler{}
	return t
}

// Current 获取当前登录用户信息
func (t *UserHandler) Current(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": shared.GlobalUser.Username,
	})
}
