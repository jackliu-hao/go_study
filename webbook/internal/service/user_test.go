package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository"
	repomocks "jikeshijian_go/webbook/internal/repository/mocks"
	"testing"
	"time"
)

func TestUserServiceV1_SingUp(t *testing.T) {
	// 定义testCase
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository
	}{
		{},
	}

	// act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 具体的测试代码
		})
	}

}

func TestUserServiceV1_Login(t *testing.T) {

	now := time.Now()

	// 定义testCase
	testCases := []struct {
		// 测试名称
		name string

		// mock
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入数据
		ctx       context.Context
		inputUser domain.User

		// 期望的数据
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 期望会调到的函数
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Id:        1,
						Email:     "123@qq.com",
						Password:  "$2a$10$by8.l433PbBlmH4Fi9WEGe7KVIk8RDber7wqjUHf1ykqvYTdg2UfS",
						CreatedAt: now,
					}, nil)

				return repo
			},

			// inputData
			ctx: context.Background(),
			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "23123123@",
			},
			//期望数据
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "$2a$10$by8.l433PbBlmH4Fi9WEGe7KVIk8RDber7wqjUHf1ykqvYTdg2UfS",
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 期望会调到的函数
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, ErrUserNotFind)
				return repo
			},
			// inputData
			ctx: context.Background(),
			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "23123123@",
			},
			//期望数据
			wantUser: domain.User{},
			wantErr:  ErrUserNotFind,
		},
		{
			name: "DB系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 期望会调到的函数
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("DB系统错误"))

				return repo
			},

			// inputData
			ctx: context.Background(),
			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "23123123@",
			},
			//期望数据
			wantUser: domain.User{},
			wantErr:  errors.New("DB系统错误"),
		},
		{
			name: "账号或密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 期望会调到的函数
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Id:        1,
						Email:     "123@qq.com",
						Password:  "$2a$10$by8.l433PbBlmH4Fi9WEGe7KVIk8RDber7wqjUHf1ykqvYTdg2UfS",
						CreatedAt: now,
					}, nil)

				return repo
			},

			// inputData
			ctx: context.Background(),
			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "aaaaa@",
			},
			//期望数据
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	// act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体的测试代码
			mockRepo := tc.mock(ctrl)
			// 保证svc不为nil
			require.NotNil(t, mockRepo)
			svc := NewUserServiceV1(mockRepo)

			//测试login
			u, err := svc.Login(tc.ctx, tc.inputUser)

			// 比较
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}
}

func TestGenPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("23123123@"), bcrypt.DefaultCost)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(hash))
}
