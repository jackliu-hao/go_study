package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"jikeshijian_go/webbook/internal/integration/startup"
	"jikeshijian_go/webbook/internal/web"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_SendSMSCode(t *testing.T) {

	rdb := startup.InitRedis()
	server := startup.InitWebServerIOC()

	testCases := []struct {
		name string
		// before
		before func(t *testing.T)
		after  func(t *testing.T)

		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功的用例",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			reqBody: `
			{"phone": "15212345678"}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		// 验证码发送太频繁
		{
			name: "验证码发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			reqBody: `
			{"phone": "15212345678"}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "验证码发送太频繁",
			},
		},
		// 非法手机号
		{
			name: "非法手机号",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			reqBody: `
			{"phone": "1578"}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 3,
				Msg:  "非法手机号",
			},
		},
		// 系统错误

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
