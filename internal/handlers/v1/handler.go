package v1

import (
	"go.bdev.cn/shortener/internal/ecodes"
	"go.bdev.cn/shortener/internal/types"
	"go.bdev.cn/shortener/internal/utils"
)

type handler struct{}

// JsonRespErr 返回错误响应
func (t *handler) JsonRespErr(errCode int) types.ResErr {
	return types.ResErr{
		ErrCode: errCode,
		ErrInfo: ecodes.GetErrCodeMessage(errCode),
	}
}

// IsURL 判断是否为URL
func (t *handler) IsURL(url string) bool {
	return utils.IsURL(url)
}
