package web

import (
	"encoding/json"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/service"
	ijwt "jikeshijian_go/webbook/internal/web/jwt"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	// 手机号
	phoneRegexPattern = "^1[3-9][0-9]{9}$"
)

type UserHandler struct {
	svc            service.UserService
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	phoneRexExp    *regexp.Regexp
	smsCodeSvc     service.CodeService
	ijwt.Handler
}

func NewUserHandler(svc service.UserService, smsSvc service.CodeService, handler ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRexExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		svc:            svc,
		smsCodeSvc:     smsSvc,
		Handler:        handler,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", h.SignUp)
	//server.PUT("/user", h.SignUp)
	//server.GET("/users/:username", h.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", h.SignUp)
	// POST /users/login
	ug.POST("/login", h.Login)

	ug.POST("/loginjwt", h.LoginJWT)
	// POST /users/edit
	ug.POST("/edit", h.Edit)
	// GET /users/profile
	ug.GET("/profile", h.Profile)
	// ProfileJWT
	ug.GET("/profileJWT", h.ProfileJWT)
	// 验证码相关路由
	ug.POST("/login_sms/code/send", h.SendSms)
	// 校验验证码
	ug.POST("/login_sms/code/check", h.verifySmsCode)
	// refreshToken
	ug.POST("/refresh_token", h.RefreshToken)
	// 退出登录
	ug.POST("/logout", h.LogoutJWT)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		// 这里其实不需要写，Bind中会自动返回错误
		//ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}
	// 调用svc
	err = h.svc.SingUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserDuplicate) {
			ctx.String(http.StatusOK, "邮箱冲突")
			return
		}
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {

	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var loginReq Req
	if err := ctx.Bind(&loginReq); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, domain.User{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	})
	if err != nil {
		if err == service.ErrInvalidUserOrPassword {
			ctx.String(http.StatusOK, "账号或密码错误")
			return
		}
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 使用jwt设置登录状态
	// 使用jwt生成一个token

	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 登录成功
	ctx.String(http.StatusOK, "登录成功")
	return

}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var loginReq Req
	if err := ctx.Bind(&loginReq); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, domain.User{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	})
	if err != nil {
		if err == service.ErrInvalidUserOrPassword {
			ctx.String(http.StatusOK, "账号或密码错误")
			return
		}
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 设置session
	sess := sessions.Default(ctx)
	// 随便设置一个session
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 30 * 60,
	})
	// 保存session
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 登录成功
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		NickName string `json:"nickName"` // 昵称
		Birthday string `json:"birthday"` // 生日
		AboutMe  string `json:"aboutMe"`  // 个人简介
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	// 获取session 的id
	session := sessions.Default(ctx)
	id := session.Get("userId")
	// 吧id转成int64
	userId, _ := id.(int64)
	err = h.svc.Edit(ctx, domain.User{
		Id:       userId,
		NickName: req.NickName,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "更新成功")
	return

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	id := session.Get("userId")
	idInt64, _ := id.(int64)
	user, err := h.svc.Profile(ctx, idInt64)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 把user转成json
	userJson, err := json.Marshal(user)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, string(userJson))

}

func (h *UserHandler) ProfileJWT(ctx *gin.Context) {

	type Profile struct {
		Id        int64     `json:"id"`
		Email     string    `json:"email"`
		NickName  string    `json:"nickName"`
		Birthday  string    `json:"birthday"`
		AboutMe   string    `json:"aboutMe"`
		Phone     string    `json:"phone"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	value, exists := ctx.Get("claims")
	if !exists {
		// 监控住这个错误
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	claims, ok := value.(*ijwt.UserClaims)
	if !ok {
		// 系统错误
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := h.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userProfile := Profile{
		Id:        user.Id,
		Email:     user.Email,
		NickName:  user.NickName,
		Birthday:  user.Birthday,
		AboutMe:   user.AboutMe,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 把user转成json
	userProfileJson, err := json.Marshal(userProfile)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, string(userProfileJson))

}

func (h *UserHandler) SendSms(context *gin.Context) {

	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	err := context.Bind(&req)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 3,
			Msg:  "非法手机号",
		})
	}
	// 校验手机号
	isPhone, err := h.phoneRexExp.MatchString(req.Phone)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isPhone {
		context.JSON(http.StatusOK, Result{
			Code: 3,
			Msg:  "非法手机号",
		})
		return
	}

	err = h.smsCodeSvc.Send(context, "login", req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrSetCodeTooManyTimes) {
			context.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "验证码发送太频繁",
			})
		} else {
			context.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
		}
		// 打印日志
		zap.L().Error("发送短信失败", zap.Error(err))
		return
	}
	context.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "发送成功",
	})
}

func (h *UserHandler) verifySmsCode(context *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	err := context.Bind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	matchString, err := h.phoneRexExp.MatchString(req.Phone)
	if err != nil {
		context.JSON(http.StatusBadRequest, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//打印日志
		zap.L().Error("验证码校验出错", zap.Error(err))
		// 日志中不能直接打印手机号码，只能在debug中打印 ，或者打印加密后的手机号 oR 脱敏
		zap.L().Debug("手机号：", zap.String("phone", req.Phone))
		return
	}
	if !matchString {
		context.JSON(http.StatusBadRequest, Result{
			Code: 5,
			Msg:  "手机号格式不对",
		})
		return
	}
	// 校验验证码是六位手机号
	verify, err := h.smsCodeSvc.Verify(context, "login", req.Phone, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrCodeVerifyFailed) {
			context.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "验证码输入错误",
			})
		} else if errors.Is(err, service.ErrCodeVerifyTooManyTimes) {
			context.JSON(http.StatusOK, Result{
				Code: 4,
				Msg:  "验证次数太多",
			})
		} else {
			context.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
		}
		return
	}
	if verify {

		// 存放ID
		domainUser, err := h.svc.FindOrCreate(context, req.Phone)
		if err != nil {
			context.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		err = h.SetLoginToken(context, domainUser.Id)
		if err != nil {
			context.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}

		context.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "发送成功",
		})
	}
	return

}

func (h *UserHandler) RefreshToken(context *gin.Context) {
	// 约定，前端在 Authorization 里面带上这个 refresh_token
	tokenStr := h.ExtractToken(context)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil || !token.Valid {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 这里如果redis崩了，可以直接跳过去

	// 需要校验此次请求的token是否有效
	err = h.CheckSession(context, rc.Ssid)
	if err != nil {
		// token 无效或者 redis 有问题
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(context, rc.Uid, rc.Ssid)
	if err != nil {
		zap.L().Error("设置Jwt异常", zap.Error(err))
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	context.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

// LogoutJWT 退出登录
func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "退出登录成功"})
}
