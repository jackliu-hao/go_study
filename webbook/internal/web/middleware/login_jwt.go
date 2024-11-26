package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"jikeshijian_go/webbook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJwtMiddlewareBuilder struct {
	paths []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
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
		tokenHeader := context.GetHeader("Authorization")
		if tokenHeader == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := segs[1]
		// 拿到cliams
		claims := &web.UserClaims{}
		parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("0776f450dd575004ba7c69930c579cae"), nil
		})
		if err != nil || !parseToken.Valid || claims.Uid <= 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 没50秒刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			signedString, err := parseToken.SignedString([]byte("0776f450dd575004ba7c69930c579cae"))
			if err != nil {
				// 记录日志
				// todo
			}
			context.Header("x-jwt-token", signedString)
		}
		// 设置jwt
		context.Set("claims", claims)

	}
}
