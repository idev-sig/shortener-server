package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener-server/internal/ecodes"
	"go.xoder.cn/shortener-server/internal/logics"
	"go.xoder.cn/shortener-server/internal/types"
)

// HistoryHandler 历史记录处理器
type HistoryHandler struct {
	handler
	logic *logics.HistoryLogic
}

// NewHistoryHandler 创建短链接处理器
func NewHistoryHandler() *HistoryHandler {
	t := &HistoryHandler{}
	t.logic = logics.NewHistoryLogic()
	return t
}

// HistoryDeleteAll 删除所有历史记录
func (t *HistoryHandler) HistoryDeleteAll(c *gin.Context) {
	var reqQuery struct {
		IDs string `form:"ids" binding:"required"`
	}
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	// log.Printf("reqQuery.IDs: %s", reqQuery.IDs)
	ids := strings.Split(reqQuery.IDs, ",")
	errCode := t.logic.HistoryDeleteAll(ids)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		c.JSON(http.StatusInternalServerError, errInfo)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HistoryList 获取历史记录列表
func (t *HistoryHandler) HistoryList(c *gin.Context) {
	var reqQuery types.ReqQueryHistory
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}
	if reqQuery.Order == "" {
		reqQuery.Order = "DESC"
	}
	if reqQuery.SortBy == "" {
		reqQuery.SortBy = "created_at"
	}

	errCode, data, pageInfo := t.logic.HistoryAll(reqQuery)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeDatabaseError {
			c.JSON(http.StatusInternalServerError, errInfo)
		} else {
			c.JSON(http.StatusBadRequest, errInfo)
		}
		return
	}

	result := types.ResSuccess[[]types.ResHistory]{
		Data: data,
		Meta: pageInfo,
	}

	c.JSON(http.StatusOK, result)
}
