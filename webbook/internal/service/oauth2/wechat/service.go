package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"jikeshijian_go/webbook/internal/domain"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://test.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx *gin.Context, state string) (string, error)
	VerifyCode(ctx *gin.Context, code string, state string) (domain.WechatInfo, error)
}

type OAuth2Service struct {
	appId     string
	appSecret string
	// 没有依赖注入
	client *http.Client
}

func (s OAuth2Service) VerifyCode(ctx *gin.Context, code string, state string) (domain.WechatInfo, error) {

	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=CODE&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	if err := decoder.Decode(&res); err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误信息，errcode:%d,errmsg:%s", res.ErrCode, res.ErrMsg)
	}
	// 登录成功
	return domain.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil

}

func NewOAuth2Service(appId string, appSecret string) Service {
	return OAuth2Service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s OAuth2Service) AuthURL(ctx *gin.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=SCOPE&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
