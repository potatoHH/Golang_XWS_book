package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("验证码发送次数太多")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrNoKnowCode             = errors.New("不知道发生什么错误,code")
)

// 编译器会在编译的时候,把set_code 的代码放进来这个luaSetCode变量中
//

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}
func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int64()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//没问题
		return nil
	case -1:
		//发送太频繁
		return ErrCodeSendTooMany
	case -2:
		//遇到这个错误,说明有人在搞
		return ErrCodeVerifyTooManyTimes
	default:
		//系统错误
		return ErrNoKnowCode
	}

}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int64()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		//没问题
		return true, nil
	case -1:
		//输入错误
		return false, nil
	case -2:
		return false, nil
	//	//系统错误
	default:
		//系统错误
		return false, errors.New("系统错误")
	}

}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone-code-%s-%s", biz, phone)
}
