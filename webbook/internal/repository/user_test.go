package repository

import (
	"context"
	"database/sql"
	"errors"
	daomocks "jikeshijian_go/webbook/internal/repository/dao/mocks"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository/cache"
	cachemocks "jikeshijian_go/webbook/internal/repository/cache/mocks"
	"jikeshijian_go/webbook/internal/repository/dao"
	"testing"
)

func TestUserRepositoryWithCache_FindById(t *testing.T) {

	//使用 time.Now().UnixMilli() 获取当前时间的毫秒级时间戳，存储在变量 nowMs 中。
	//使用 time.UnixMilli(nowMs) 将毫秒级时间戳转换为 time.Time 类型的时间对象，存储在变量 now 中。
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)

	testCases := []struct {
		// 测试用例名称
		name string
		//mock
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		//input
		userId int64
		//expect
		wantUser domain.User
		wantErr  error
	}{
		// 缓存未命中
		{
			name: "查询成功，但是缓存结果不存在啊",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(1)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)
				// 还调了set方法
				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:        uid,
					Email:     "123@qq.com",
					Password:  "123456",
					Birthday:  "100",
					AboutMe:   "自我介绍",
					Phone:     "15212345678",
					CreatedAt: now,
					UpdatedAt: time.UnixMilli(102),
				}).
					Return(nil)

				userDao.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Birthday: "100",
						AboutMe:  "自我介绍",
						Phone: sql.NullString{
							String: "15212345678",
							Valid:  true,
						},
						CreatedAt: nowMs,
						UpdatedAt: 102,
					}, nil)

				return userDao, userCache
			},
			userId: 1,
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "123456",
				Birthday:  "100",
				AboutMe:   "自我介绍",
				Phone:     "15212345678",
				CreatedAt: now,
				UpdatedAt: time.UnixMilli(102),
			},
			wantErr: nil,
		},
		// 缓存命中
		{
			name: "查询成功，缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(1)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{
						Id:        uid,
						Email:     "123@qq.com",
						Password:  "123456",
						Birthday:  "100",
						AboutMe:   "自我介绍",
						Phone:     "15212345678",
						CreatedAt: now,
						UpdatedAt: time.UnixMilli(102),
					}, nil)

				return userDao, userCache
			},
			userId: 1,
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "123456",
				Birthday:  "100",
				AboutMe:   "自我介绍",
				Phone:     "15212345678",
				CreatedAt: now,
				UpdatedAt: time.UnixMilli(102),
			},
			wantErr: nil,
		},
		// 查询失败
		{
			name: "查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(1)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)

				userDao.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{}, errors.New("查询失败"))

				return userDao, userCache
			},
			userId:   1,
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		// 缓存设置失败
		{
			name: "查询成功，但是缓存失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(1)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrKeyNotExist)
				// 还调了set方法
				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:        uid,
					Email:     "123@qq.com",
					Password:  "123456",
					Birthday:  "100",
					AboutMe:   "自我介绍",
					Phone:     "15212345678",
					CreatedAt: now,
					UpdatedAt: time.UnixMilli(102),
				}).
					Return(errors.New("设置缓存失败"))

				userDao.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Birthday: "100",
						AboutMe:  "自我介绍",
						Phone: sql.NullString{
							String: "15212345678",
							Valid:  true,
						},
						CreatedAt: nowMs,
						UpdatedAt: 102,
					}, nil)

				return userDao, userCache
			},
			userId: 1,
			wantUser: domain.User{
				Id:        1,
				Email:     "123@qq.com",
				Password:  "123456",
				Birthday:  "100",
				AboutMe:   "自我介绍",
				Phone:     "15212345678",
				CreatedAt: now,
				UpdatedAt: time.UnixMilli(102),
			},
			wantErr: errors.New("设置缓存失败"),
		},
		// 查询缓存出错
		{
			name: "查询缓存出错",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				uid := int64(1)
				userDao := daomocks.NewMockUserDao(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)

				userCache.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, errors.New("查询缓存出错"))

				return userDao, userCache
			},
			userId:   1,
			wantUser: domain.User{},
			wantErr:  errors.New("查询缓存出错"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userDao, userCache := tc.mock(ctrl)
			repo := NewUserRepositoryWithCache(userDao, userCache)
			user, err := repo.FindById(context.Background(), tc.userId)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}

}
