package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/service"
	svcmocks "jikeshijian_go/webbook/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	claims := c.(*UserClaims)
	fmt.Println(claims.Uid)
}

func TestUserHandler_SignUp(t *testing.T) {

	testCases := []struct {
		// 测试名称
		name string
		// mock service
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 输入用例
		body string

		// 期望的返回Code
		wantCode int
		// 期望的返回值
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), gomock.Any()).
					Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123@qq.com","password": "1234Aa1!56","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数错误 ， Bind失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": 123@qq.com","password": "1234Aa1!56","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusBadRequest,
			// 不需要返回
			//wantBody: "注册成功",
		},
		{
			name: "非法邮箱格式",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123qq.com","password": "1234Aa1!56","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusOK,
			// 不需要返回
			wantBody: "非法邮箱格式",
		},
		{
			name: "两次密码输入不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123@qq.com","password": "1234Aa1!056","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusOK,
			wantBody: "两次输入密码不对",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123@qq.com","password": "123456","confirmPassword": "123456"}`,
			wantCode: http.StatusOK,
			wantBody: "密码必须包含字母、数字、特殊字符，并且不少于八位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "1234Aa1!56",
				}).
					Return(service.ErrUserDuplicate)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123@qq.com","password": "1234Aa1!56","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SingUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "1234Aa1!56",
				}).
					Return(errors.New("系统错误"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			body:     `{"email": "123@qq.com","password": "1234Aa1!56","confirmPassword": "1234Aa1!56"}`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 注册路由
			server := gin.Default()
			// 创建一个控制器
			ctrl := gomock.NewController(t)
			// 结束后，需要会在Finnish比较
			defer ctrl.Finish()
			userService, codeService := tc.mock(ctrl)
			NewUserHandler(userService, codeService).RegisterRoutes(server)
			// mock
			req, err := http.NewRequest(http.MethodPost,
				"/users/signup",
				bytes.NewBuffer([]byte(tc.body)))
			// 保证不存在error
			require.NoError(t, err)
			// 继续使用req
			req.Header.Set("Content-Type", "application/json")
			// 构造响应
			resp := httptest.NewRecorder()

			// 开启server,Gin框架的入口
			// 这样调用的时候gin会执行这个请求，响应写到resp
			server.ServeHTTP(resp, req)

			// 判断是否成功或失败
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}

	time.Sleep(3 * time.Second)

}

func TestMock(t *testing.T) {
	// 创建一个控制器
	ctrl := gomock.NewController(t)
	// 结束后，需要会在Finnish比较
	defer ctrl.Finish()

	// 模拟service
	userSvc := svcmocks.NewMockUserService(ctrl)
	// 发起期望调用 ，因为singUp返回是一个error，所以在调用的时候也需要返回一个error
	userSvc.EXPECT().SingUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	// 发起实际调用
	// context.Background() 从创建一个空白的上下文
	err := userSvc.SingUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})

	// 这里返回的就是 mock error
	t.Log(err)

}
