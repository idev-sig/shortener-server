package v1

import (
	"go.xoder.cn/shortener-server/internal/ecodes"
	"go.xoder.cn/shortener-server/internal/types"
	"go.xoder.cn/shortener-server/internal/utils"
)

type handler struct{}

// JsonRespErr 返回错误响应 (符合 OpenAPI ErrorResponse)
func (t *handler) JsonRespErr(errCode int) types.ResErr {
	return types.ResErr{
		ErrorCode:    utils.IntToString(errCode),
		ErrorMessage: ecodes.GetErrCodeMessage(errCode),
	}
}

// IsURL 判断是否为URL
func (t *handler) IsURL(url string) bool {
	return utils.IsURL(url)
}
