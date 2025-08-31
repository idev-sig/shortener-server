package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener/internal/ecodes"
	"go.xoder.cn/shortener/internal/logics"
)

// AccountHandler 账号处理器
type AccountHandler struct {
	handler
	logic *logics.AccountLogic
}

// NewAccountHandler 创建账号处理器
func NewAccountHandler() *AccountHandler {
	t := &AccountHandler{}
	return t
}

// Login 账号登录
func (t *AccountHandler) Login(c *gin.Context) {
	var reqJson struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Auto     bool   `json:"auto,omitempty"`
	}

	if err := c.ShouldBindJSON(&reqJson); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	token, err := t.logic.Login(reqJson.Username, reqJson.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeUserPasswordError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Logout 账号登出
func (t *AccountHandler) Logout(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	_ = t.logic.Remove(token)

	c.JSON(http.StatusNoContent, nil)
}
