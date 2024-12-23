package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"jikeshijian_go/webbook/internal/domain"
	"jikeshijian_go/webbook/internal/repository/cache"
	"jikeshijian_go/webbook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache cache.UserCache
}

// Create 创建
func (r *UserRepository) Create(ctx context.Context, user domain.User) error {

	// 1. 密码加密

	// 2. 插入到数据库

	// 3. 返回结果
	return r.dao.Insert(ctx, r.domain2Entity(user))

}

func NewUserRepository(dao *dao.UserDAO, userCache cache.UserCache) UserRepository {

	return UserRepository{
		dao:   dao,
		cache: userCache,
	}
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	// 先从cache中查找
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 数据存在
		return u, nil
	}
	// 数据不存在
	if err == cache.ErrKeyNotExist {
		// 从数据库中加载
		ud, err := r.dao.FindById(ctx, id)
		if err != nil {
			return domain.User{}, ErrUserNotFound
		}
		user := r.entity2Domain(ud)
		go func() {
			// 找到了在写到cache
			err = r.cache.Set(ctx, user)
			if err != nil {
				// how to do ?
				// 应该打日志，做监控
			}
		}()
		return user, err
	}
	// 这里怎么搞？缓存出错了
	return domain.User{}, err

}

func (r *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	// 需要将PO转成BO ,返回给service层
	if err != nil {
		return domain.User{}, err
	}
	return r.entity2Domain(u), nil
}

func (r *UserRepository) UpdateById(ctx *gin.Context, user domain.User) error {
	err := r.dao.UpdateById(ctx, dao.User{
		Id:       user.Id,
		NickName: user.NickName,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday,
	})
	return err
}

func (r *UserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	// 需要将PO转成BO ,返回给service层
	if err != nil {
		return domain.User{}, err
	}
	return r.entity2Domain(u), nil
}

func (r *UserRepository) entity2Domain(ud dao.User) domain.User {

	return domain.User{
		Id:        ud.Id,
		Email:     ud.Email.String,
		Password:  ud.Password,
		CreatedAt: time.UnixMilli(ud.CreatedAt),
		UpdatedAt: time.UnixMilli(ud.UpdatedAt),
		AboutMe:   ud.AboutMe,
		Birthday:  ud.Birthday,
		NickName:  ud.NickName,
		Phone:     ud.Phone.String,
	}
}

func (r *UserRepository) domain2Entity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password:  u.Password,
		CreatedAt: u.CreatedAt.UnixMilli(),
		UpdatedAt: u.UpdatedAt.UnixMilli(),
		AboutMe:   u.AboutMe,
		Birthday:  u.Birthday,
		NickName:  u.NickName,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
	}
}
