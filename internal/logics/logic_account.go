package logics

import (
	"errors"

	"go.xoder.cn/shortener/internal/shared"
	"go.xoder.cn/shortener/internal/utils"
)

// AccountLogic 账号处理逻辑
type AccountLogic struct {
	logic
}

// NewAccountLogic 账号逻辑实例化
func NewAccountLogic() *AccountLogic {
	t := &AccountLogic{}
	t.init()
	return t
}

// Check 检查账号密码
func (t *AccountLogic) Login(username string, password string) (string, error) {
	if username != shared.GlobalUser.Username || password != shared.GlobalUser.Password {
		return "", errors.New("用户名或密码错误")
	}

	userToken := utils.GenerateCode(32)

	shared.GlobalUserCache.Range(func(key, value any) bool {
		shared.GlobalUserCache.Delete(key)
		return true
	})
	shared.GlobalUserCache.Store(userToken, username)

	return userToken, nil
}

// Remove 移除账号
func (t *AccountLogic) Remove(token string) error {
	shared.GlobalUserCache.Delete(token)
	return nil
}
