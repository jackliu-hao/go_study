package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	ijwt "jikeshijian_go/webbook/internal/web/jwt"
	"net/http"
)

type LoginJwtMiddlewareBuilder struct {
	paths      []string
	jwtHandler ijwt.Handler
	redis      redis.Cmdable
}

func NewLoginJwtMiddlewareBuilder(handler ijwt.Handler) *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{
		jwtHandler: handler,
	}
}

func (l *LoginJwtMiddlewareBuilder) IgnorePath(path string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {

	return func(context *gin.Context) {

		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}
		//if context.Request.URL.Path == "/users/login" ||
		//	context.Request.URL.Path == "/users/signup" {
		//	return
		//}

		// 校验是否登录
		// 使用jwt校验
		//tokenHeader := context.GetHeader("Authorization")
		//if tokenHeader == "" {
		//	context.AbortWithStatus(http.StatusUnauthorized)
		//}
		//segs := strings.Split(tokenHeader, " ")
		//if len(segs) != 2 {
		//	context.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//token := segs[1]
		// 拿到cliams
		tokenStr := l.jwtHandler.ExtractToken(context)
		var uc ijwt.UserClaims
		parseToken, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})

		if err != nil || !parseToken.Valid || uc.Uid <= 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 这里如果redis崩了，可以直接跳过去

		// 需要校验此次请求的token是否有效
		// 这里看
		err = l.jwtHandler.CheckSession(context, uc.Ssid)
		if err != nil {
			// token 无效或者 redis 有问题
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 没50秒刷新一次
		// 引入长短token后，这个就不需要了
		//now := time.Now()
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	signedString, err := parseToken.SignedString([]byte("0776f450dd575004ba7c69930c579cae"))
		//	if err != nil {
		//		// 记录日志
		//		//
		//	}
		//	context.Header("x-jwt-token", signedString)
		//}
		// 设置jwt
		context.Set("claims", uc)

	}
}
