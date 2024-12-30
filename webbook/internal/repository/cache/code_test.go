package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	redismocks "jikeshijian_go/webbook/mocks"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		// mock
		mock func(ctrl *gomock.Controller) redis.Cmdable

		// input
		biz   string
		phone string
		code  string

		//except
		wantErr error
	}{
		// 验证码发送成功
		{
			name: "验证码发送成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rdb := redismocks.NewMockCmdable(ctrl)

				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))

				rdb.EXPECT().Eval(context.Background(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "12345678901")},
					"123456",
				).
					Return(cmd)

				return rdb
			},
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: nil,
		},
		// 发送太频繁
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rdb := redismocks.NewMockCmdable(ctrl)

				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))

				rdb.EXPECT().Eval(context.Background(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "12345678901")},
					"123456",
				).
					Return(cmd)

				return rdb
			},
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: ErrSetCodeTooManyTimes,
		},
		//系统错误
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rdb := redismocks.NewMockCmdable(ctrl)

				cmd := redis.NewCmd(context.Background())
				//cmd.SetErr(ErrUnkonwForCode)
				cmd.SetVal(int64(5))

				rdb.EXPECT().Eval(context.Background(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "12345678901")},
					"123456",
				).
					Return(cmd)

				return rdb
			},
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: ErrUnkonwForCode,
		},
		// 脚本执行出错
		{
			name: "脚本执行出错",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rdb := redismocks.NewMockCmdable(ctrl)

				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("脚本执行出错"))
				cmd.SetVal(int64(5))

				rdb.EXPECT().Eval(context.Background(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "12345678901")},
					"123456",
				).
					Return(cmd)

				return rdb
			},
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: errors.New("脚本执行出错"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisCmdable := tc.mock(ctrl)
			cache := NewRedisCodeCache(redisCmdable)
			err := cache.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)

		})
	}

}
