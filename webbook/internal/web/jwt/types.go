package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ClearToken(ctx *gin.Context) error
	ExtractToken(ctx *gin.Context) string
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	CheckSession(ctx *gin.Context, ssid string) error
}

// RefreshClaims 刷新cookie
type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
	// 用于标识是否退出登录
	Ssid string
}

// UserClaims 存放jwt的内容
type UserClaims struct {
	jwt.RegisteredClaims
	// 声名自己要放到Claim的数据
	Uid int64
	// 用于标识是否退出登录
	Ssid string
	// User-Agent

}
