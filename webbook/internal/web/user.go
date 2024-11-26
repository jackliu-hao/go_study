package web

import (
	"encoding/json"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/service"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	svc            *service.UserService
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
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
	ug.GET("/profileJWT", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
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
		if err == service.ErrDuplicateEmail {
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

	claims := UserClaims{
		// 设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		},
		// 设置用户id
		Uid: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("0776f450dd575004ba7c69930c579cae"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	// 把jwt放到header中
	ctx.Header("x-jwt-token", signedString)
	fmt.Println(signedString)

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
	value, exists := ctx.Get("claims")
	if !exists {
		// 监控住这个错误
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	claims, ok := value.(*UserClaims)
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

	// 把user转成json
	userJson, err := json.Marshal(user)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, string(userJson))

}

// UserClaims 存放jwt的内容
type UserClaims struct {
	jwt.RegisteredClaims
	// 声名自己要放到Claim的数据
	Uid int64
}
