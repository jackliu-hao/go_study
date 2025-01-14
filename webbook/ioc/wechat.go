package ioc

import (
	"jikeshijian_go/webbook/internal/service/oauth2/wechat"
	"jikeshijian_go/webbook/internal/web"
)

func InitOAuth2WechatService() wechat.Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	appID, ok := "123", true
	if !ok {
		panic("找不到环境变量 WECHAT_APP_ID")
	}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	appSecret, ok := "123", true
	if !ok {
		panic("找不到环境变量 WECHAT_APP_SECRET")
	}
	return wechat.NewOAuth2Service(appID, appSecret)
}

func NewWechatHandler() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
