package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {

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

		//// 校验是否登录
		session := sessions.Default(context)
		//if session == nil {
		//	context.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		id := session.Get("userId")
		if id == nil {
			// 不存在 id
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 已经登录了，刷新session
		// 拿到上一次的更新时间
		updateTime := session.Get("update_time")
		now := time.Now().UnixMilli()
		// 还没有刷新过
		if updateTime == nil {
			// 刚登陆，还没刷新
			session.Set("update_time", now)
			session.Save()
			return
		}
		// 已经刷新过了
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			// 系统错误
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 60*1000 {
			// 刚登陆，还没刷新
			session.Set("update_time", now)
			session.Save()
			return
		}

	}
}
