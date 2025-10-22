package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go.xoder.cn/shortener/internal/ecodes"
	"go.xoder.cn/shortener/internal/logics"
	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/types"
	"go.xoder.cn/shortener/internal/utils"
)

// ShortenHandler 短链接处理器
type ShortenHandler struct {
	handler
	logic *logics.ShortenLogic
}

// NewShortenHandler 创建短链接处理器
func NewShortenHandler() *ShortenHandler {
	t := &ShortenHandler{}
	t.logic = logics.NewShortenLogic()
	return t
}

// ShortenRedirect 短链接跳转
func (t *ShortenHandler) ShortenRedirect(c *gin.Context) {
	var reqUri types.ReqCode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	errCode, data := t.logic.ShortenFind(reqUri.Code)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeNotFound {
			c.JSON(http.StatusNotFound, errInfo)
		} else {
			c.JSON(http.StatusInternalServerError, errInfo)
		}
		return
	}

	// 异步记录访问历史
	record := logics.NewHistoryLogic()
	go func() {
		_ = record.HistoryAdd(
			types.HistoryParams{
				URLID:     data.ID,
				ShortCode: data.ShortCode,
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Referer:   c.Request.Referer(),
			},
		)
	}()

	c.Redirect(http.StatusFound, data.OriginalURL)
}

// ShortenAdd 添加短链接
func (t *ShortenHandler) ShortenAdd(c *gin.Context) {
	var reqJson struct {
		ShortCode   string `json:"short_code,omitempty"`
		OriginalURL string `json:"original_url" binding:"required,url"`
		Description string `json:"description,omitempty"`
	}

	if err := c.ShouldBindJSON(&reqJson); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	if reqJson.OriginalURL != "" && !t.IsURL(reqJson.OriginalURL) {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	// 生成短码
	if reqJson.ShortCode == "" {
		reqJson.ShortCode = utils.GenerateCode(shared.GlobalShorten.Length)
	}

	if len(reqJson.ShortCode) > 16 {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeBadRequest))
		return
	}

	errCode, data := t.logic.ShortenAdd(reqJson.ShortCode, reqJson.OriginalURL, reqJson.Description)
	if errCode != 0 {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeConflict {
			c.JSON(http.StatusConflict, errInfo)
		} else {
			c.JSON(http.StatusInternalServerError, errInfo)
		}
		return
	}

	c.Header("Location", c.Request.RequestURI+"/"+data.ShortCode)
	c.JSON(http.StatusCreated, data)
}

// ShortenDelete 删除短链接
func (t *ShortenHandler) ShortenDelete(c *gin.Context) {
	var reqUri types.ReqCode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	errCode := t.logic.ShortenDelete(reqUri.Code)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeNotFound {
			c.JSON(http.StatusNotFound, errInfo)
		} else {
			c.JSON(http.StatusInternalServerError, errInfo)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ShortenDeleteAll 删除所有短链接
func (t *ShortenHandler) ShortenDeleteAll(c *gin.Context) {
	var reqQuery struct {
		IDs string `form:"ids" binding:"required"`
	}
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	// log.Printf("reqQuery.IDs: %s", reqQuery.IDs)
	ids := strings.Split(reqQuery.IDs, ",")
	errCode := t.logic.ShortenDeleteAll(ids)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		c.JSON(http.StatusInternalServerError, errInfo)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ShortenUpdate 更新短链接
func (t *ShortenHandler) ShortenUpdate(c *gin.Context) {
	var reqUri types.ReqCode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	var reqJson struct {
		OriginalURL string `json:"original_url,omitempty" binding:"omitempty,url"`
		Description string `json:"description,omitempty"`
	}
	if err := c.ShouldBindJSON(&reqJson); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	if reqJson.OriginalURL != "" && !t.IsURL(reqJson.OriginalURL) {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	errCode, data := t.logic.ShortenUpdate(reqUri.Code, reqJson.OriginalURL, reqJson.Description)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeNotFound {
			c.JSON(http.StatusNotFound, errInfo)
		} else {
			c.JSON(http.StatusInternalServerError, errInfo)
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// ShortenFind 获取短链接
func (t *ShortenHandler) ShortenFind(c *gin.Context) {
	var reqUri types.ReqCode
	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}

	errCode, data := t.logic.ShortenFind(reqUri.Code)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeNotFound {
			c.JSON(http.StatusNotFound, errInfo)
		} else {
			c.JSON(http.StatusInternalServerError, errInfo)
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// ShortenList 获取短链接列表
func (t *ShortenHandler) ShortenList(c *gin.Context) {
	var reqQuery types.ReqQueryShorten
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		// log.Printf("err: %v", err)
		c.JSON(http.StatusBadRequest, t.JsonRespErr(ecodes.ErrCodeInvalidParam))
		return
	}
	if reqQuery.Order == "" {
		reqQuery.Order = "DESC"
	}
	if reqQuery.SortBy == "" {
		reqQuery.SortBy = "created_at"
	}

	errCode, data, pageInfo := t.logic.ShortenAll(reqQuery)
	if errCode != ecodes.ErrCodeSuccess {
		errInfo := t.JsonRespErr(errCode)
		if errCode == ecodes.ErrCodeDatabaseError {
			c.JSON(http.StatusInternalServerError, errInfo)
		} else {
			c.JSON(http.StatusBadRequest, errInfo)
		}
		return
	}

	result := types.ResSuccess[[]types.ResShorten]{
		Data: data,
		Meta: pageInfo,
	}

	c.JSON(http.StatusOK, result)
}
