package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"jikeshijian_go/webbook/internal/service/sms"
)

type AuthSMSService struct {
	svc sms.Service
	key string
}

// Send 发送，其中biz 必须是线下申请的代表业务方的jwt
func (a AuthSMSService) Send(ctx context.Context, biz string, args []string, phoneNumbers ...string) error {

	var tc TokenClaims
	// 权限校验
	//第一个参数：biz
	//类型：string 或 []byte
	//描述：这是待解析的JWT字符串或字节数组。通常是从请求头或其他地方获取的未解码的JWT。
	//第二个参数：&tc
	//类型：指向自定义声明结构体的指针（例如 *CustomClaims）
	//描述：这是一个指向自定义声明结构体的指针，用于存储解析后的JWT声明信息。你需要确保这个结构体实现了 jwt.Claims 接口。
	//第三个参数：密钥回调函数
	//类型：func(*jwt.Token) (interface{}, error)
	//描述：这是一个回调函数，用于提供签名验证所需的密钥。该函数接收一个 *jwt.Token 作为参数，并返回一个接口类型的密钥和一个错误。
	//在这个例子中，回调函数返回了 a.key 和 nil，表示没有错误发生。
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		// 如果报错，说明对应的业务方存在问题
		return err
	}
	if !token.Valid {
		// 如果校验失败，说明对应的业务方存在问题
		return errors.New("token 不合法")
	}

	return a.svc.Send(ctx, tc.Tpl, args, phoneNumbers...)
}

type TokenClaims struct {
	Tpl string
	jwt.RegisteredClaims
}
