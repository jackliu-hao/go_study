package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"jikeshijian_go/webbook/internal/web"
	ijwt "jikeshijian_go/webbook/internal/web/jwt"
	"jikeshijian_go/webbook/internal/web/middleware"
	"jikeshijian_go/webbook/pkg/ginx/middlewares/ratelimit"
	"jikeshijian_go/webbook/pkg/logger"
	ratelimit2 "jikeshijian_go/webbook/pkg/ratelimit"
	"strings"
	"time"
)

//func InitWebServerV1(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
//	server := gin.Default()
//	server.Use(mdls...)
//	for _, hdl := range hdls {
//		hdl.RegisterRoutes(server)
//	}
//	//userHdl.RegisterRoutes(server)
//	return server
//}

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oath2Wechathdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	oath2Wechathdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, jwtHandler ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {

	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,

			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		// logger
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
			l.Debug("HTTP请求", logger.Field{Key: "al", Val: al})
		}).AllowRespBody().AllowReqBody().Build(),
		func(ctx *gin.Context) {
			println("这是我的 Middleware")
		},
		ratelimit.NewBuilder(ratelimit2.NewRedisSlidingWindowLimiter(redisClient, 10*time.Second, 100)).Build(),
		middleware.NewLoginJwtMiddlewareBuilder(jwtHandler).
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			IgnorePath("/users/login").
			IgnorePath("/users/register").
			IgnorePath("/users/loginjwt").
			IgnorePath("/oauth2/wechat/authurl").
			IgnorePath("/oauth2/wechat/callback").
			IgnorePath("/users/refresh_token").
			Build(),
	}
}
