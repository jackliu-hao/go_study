package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"jikeshijian_go/webbook/internal/service"
	"jikeshijian_go/webbook/internal/service/oauth2/wechat"
	ijwt "jikeshijian_go/webbook/internal/web/jwt"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc        wechat.Service
	userSvc    service.UserService
	jwtHandler ijwt.Handler
	cfg        WechatHandlerConfig
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService, cfg WechatHandlerConfig, jwtHandler ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:        svc,
		userSvc:    userSvc,
		cfg:        cfg,
		jwtHandler: jwtHandler,
	}
}

type WechatHandlerConfig struct {
	Secure bool
}

func (h OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {

	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {

	state := uuid.New()

	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "构造扫码登录url失败",
		})
	}
	err = h.setCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setCookie(ctx *gin.Context, state string) error {
	// 传一下 cookie
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间，预期三分钟
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})
	signedString, err := token.SignedString(ijwt.JWTKey)
	if err != nil {
		return fmt.Errorf("生成jwt失败, %w", err)
	}
	ctx.SetCookie("jwt-state", signedString, 600, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

func (h OAuth2WechatHandler) Callback(context *gin.Context) {
	code := context.Query("code")
	state := context.Query("state")
	err := h.verifyState(context, state)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}

	wechatInfo, err := h.svc.VerifyCode(context, code, state)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "登录失败",
		})
		return
	}
	// 从 user_service中拿到uid
	weChatUser, err := h.userSvc.FindOrCreateByWechat(context, wechatInfo)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "系统错误",
		})
		return
	}
	// 这里做什么
	err = h.jwtHandler.SetLoginToken(context, weChatUser.Id)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}
	context.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (h OAuth2WechatHandler) verifyState(context *gin.Context, state string) error {

	// 先校验state是否相同
	cookie, err := context.Cookie("jwt-state")
	if err != nil {
		// 疑似攻击
		return fmt.Errorf("拿不到state的cookie, %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(cookie, &sc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.JWTKey, nil
	})

	if err != nil || !token.Valid {

		return fmt.Errorf("验证码过期，刷新重试 , %w ", err)
	}
	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}
