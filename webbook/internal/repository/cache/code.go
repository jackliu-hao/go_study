package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSetCodeTooManyTimes    = fmt.Errorf("验证码发送次数过多")
	ErrCodeVerifyFailed       = fmt.Errorf("验证码错误")
	ErrCodeVerifyTooManyTimes = fmt.Errorf("验证码验证次数过多")
	ErrUnkonwForCode          = fmt.Errorf("未知错误")
)

// 这段Go代码使用了go:embed指令，将名为lua/set_code.lua的Lua脚本文件嵌入到Go二进制文件中，并将其内容存储在变量luaSetCode中。
// luaSetCode：存储嵌入文件内容的字符串变量。
//
//go:embed：用于将指定文件的内容嵌入到Go程序中。
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func (c *CodeCache) Set(ctx context.Context, biz string,
	phone string, code string) error {
	// 使用Lua脚本通过Eval方法执行设置操作，该操作在Redis中实现。
	// 该方法主要用于设置验证码，其中包含的参数如下：
	// - ctx: 上下文，用于传递请求和超时信息。
	// - luaSetCode: Lua脚本，定义了如何在Redis中设置和过期验证码。
	// - key(biz, phone): 生成Redis键的函数，基于业务标识和电话号码创建唯一的键。
	// - code: 需要存储的验证码。
	// 该方法返回两个值：
	// - res: 从Eval方法返回的结果，通常表示操作状态。
	// - err: 可能发生的错误，如果执行成功，则为nil。
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	// 判断下lua脚本返回的结果
	switch res {
	case 0:
		// 验证码发送成功
		return nil
	case -1:
		// 发送太频繁
		return ErrSetCodeTooManyTimes
	default:
		// 系统错误
		return ErrUnkonwForCode
	}

}

func (c *CodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		// redis 坏了
		return false, err
	}
	switch res {
	case 0:
		// 验证成功
		return true, nil
	case -1:
		// 验证码次数太多了
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		// 验证码错误
		return false, ErrCodeVerifyFailed
	default:
		// 未知错误
		return false, ErrUnkonwForCode
	}
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)

}
